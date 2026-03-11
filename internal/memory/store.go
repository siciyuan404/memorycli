package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/memorycli/memorycli/internal/redis"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Store(ctx context.Context, memory *Memory) error {
	data, err := json.Marshal(memory)
	if err != nil {
		return fmt.Errorf("failed to marshal memory: %w", err)
	}

	key := fmt.Sprintf("memory:%s", memory.ID)
	if err := s.client.Set(ctx, key, data, 0); err != nil {
		return fmt.Errorf("failed to store memory: %w", err)
	}

	agentKey := fmt.Sprintf("agent:%s:memories", memory.AgentID)
	if err := s.client.SAdd(ctx, agentKey, memory.ID); err != nil {
		return fmt.Errorf("failed to add memory to agent index: %w", err)
	}

	typeKey := fmt.Sprintf("memory:type:%s", memory.Type)
	if err := s.client.SAdd(ctx, typeKey, memory.ID); err != nil {
		return fmt.Errorf("failed to add memory to type index: %w", err)
	}

	for _, tag := range memory.Tags {
		tagKey := fmt.Sprintf("memory:tag:%s", tag)
		if err := s.client.SAdd(ctx, tagKey, memory.ID); err != nil {
			return fmt.Errorf("failed to add memory to tag index: %w", err)
		}
	}

	if memory.IsShared() {
		if err := s.client.SAdd(ctx, "memory:shared", memory.ID); err != nil {
			return fmt.Errorf("failed to add memory to shared index: %w", err)
		}
	}

	if memory.TTL > 0 {
		if err := s.client.Expire(ctx, key, time.Duration(memory.TTL)*time.Second); err != nil {
			return fmt.Errorf("failed to set memory expiration: %w", err)
		}
	}

	return nil
}

func (s *Store) Get(ctx context.Context, id string) (*Memory, error) {
	key := fmt.Sprintf("memory:%s", id)
	data, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	var memory Memory
	if err := json.Unmarshal([]byte(data), &memory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memory: %w", err)
	}

	return &memory, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	memory, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("memory:%s", id)
	if err := s.client.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	agentKey := fmt.Sprintf("agent:%s:memories", memory.AgentID)
	if err := s.client.SRem(ctx, agentKey, id); err != nil {
		return fmt.Errorf("failed to remove memory from agent index: %w", err)
	}

	typeKey := fmt.Sprintf("memory:type:%s", memory.Type)
	if err := s.client.SRem(ctx, typeKey, id); err != nil {
		return fmt.Errorf("failed to remove memory from type index: %w", err)
	}

	for _, tag := range memory.Tags {
		tagKey := fmt.Sprintf("memory:tag:%s", tag)
		if err := s.client.SRem(ctx, tagKey, id); err != nil {
			return fmt.Errorf("failed to remove memory from tag index: %w", err)
		}
	}

	if memory.IsShared() {
		if err := s.client.SRem(ctx, "memory:shared", id); err != nil {
			return fmt.Errorf("failed to remove memory from shared index: %w", err)
		}
	}

	return nil
}

func (s *Store) Share(ctx context.Context, id string, targetAgents []string) error {
	memory, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	for _, agentID := range targetAgents {
		memory.ShareWith(agentID)
	}

	memory.UpdatedAt = time.Now()

	if !memory.IsShared() {
		if err := s.client.SAdd(ctx, "memory:shared", id); err != nil {
			return fmt.Errorf("failed to add memory to shared index: %w", err)
		}
	}

	return s.Store(ctx, memory)
}

func (s *Store) Search(ctx context.Context, query string, opts SearchOptions) ([]*Memory, error) {
	var memoryIDs []string
	var err error

	if opts.AgentID != "" {
		agentKey := fmt.Sprintf("agent:%s:memories", opts.AgentID)
		memoryIDs, err = s.client.SMembers(ctx, agentKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent memories: %w", err)
		}
	}

	if opts.MemoryType != "" {
		typeKey := fmt.Sprintf("memory:type:%s", opts.MemoryType)
		typeIDs, err := s.client.SMembers(ctx, typeKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get type memories: %w", err)
		}

		if len(memoryIDs) > 0 {
			memoryIDs = intersect(memoryIDs, typeIDs)
		} else {
			memoryIDs = typeIDs
		}
	}

	if opts.SharedOnly {
		sharedIDs, err := s.client.SMembers(ctx, "memory:shared")
		if err != nil {
			return nil, fmt.Errorf("failed to get shared memories: %w", err)
		}

		if len(memoryIDs) > 0 {
			memoryIDs = intersect(memoryIDs, sharedIDs)
		} else {
			memoryIDs = sharedIDs
		}
	}

	if len(memoryIDs) == 0 {
		pattern := "memory:*"
		allKeys, err := s.client.Keys(ctx, pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to get all memories: %w", err)
		}

		memoryIDs = make([]string, 0, len(allKeys))
		for _, key := range allKeys {
			memoryIDs = append(memoryIDs, strings.TrimPrefix(key, "memory:"))
		}
	}

	memories := make([]*Memory, 0)
	queryLower := strings.ToLower(query)

	for _, id := range memoryIDs {
		memory, err := s.Get(ctx, id)
		if err != nil {
			continue
		}

		if query != "" && !strings.Contains(strings.ToLower(memory.Content), queryLower) {
			continue
		}

		memories = append(memories, memory)
	}

	if len(memories) > opts.Limit {
		memories = memories[:opts.Limit]
	}

	return memories, nil
}

func (s *Store) List(ctx context.Context, opts SearchOptions) ([]*Memory, error) {
	return s.Search(ctx, "", opts)
}

type SearchOptions struct {
	AgentID     string
	MemoryType  MemoryType
	SharedOnly  bool
	Limit       int
}

func intersect(a, b []string) []string {
	m := make(map[string]bool)
	for _, item := range a {
		m[item] = true
	}

	var result []string
	for _, item := range b {
		if m[item] {
			result = append(result, item)
		}
	}

	return result
}

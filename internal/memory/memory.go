package memory

import (
	"time"

	"github.com/google/uuid"
)

type MemoryType string

const (
	MemoryTypeKnowledge    MemoryType = "knowledge"
	MemoryTypePreference   MemoryType = "preference"
	MemoryTypeConversation MemoryType = "conversation"
	MemoryTypeTask         MemoryType = "task"
)

type Memory struct {
	ID         string         `json:"id"`
	AgentID    string         `json:"agent_id"`
	Type       MemoryType     `json:"type"`
	Content    string         `json:"content"`
	Embedding  []float64      `json:"embedding,omitempty"`
	Importance float64        `json:"importance"`
	Tags       []string       `json:"tags"`
	SharedWith []string       `json:"shared_with,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	TTL        int            `json:"ttl,omitempty"`
}

func NewMemory(content string, memoryType MemoryType, agentID string) *Memory {
	now := time.Now()
	return &Memory{
		ID:         generateID(),
		AgentID:    agentID,
		Type:       memoryType,
		Content:    content,
		Importance: 0.5,
		Tags:       []string{},
		SharedWith: []string{},
		Metadata:   make(map[string]any),
		CreatedAt:  now,
		UpdatedAt:  now,
		TTL:        0,
	}
}

func generateID() string {
	return "mem_" + uuid.New().String()[:8]
}

func (m *Memory) SetImportance(importance float64) {
	if importance < 0 {
		importance = 0
	} else if importance > 1 {
		importance = 1
	}
	m.Importance = importance
}

func (m *Memory) AddTag(tag string) {
	for _, t := range m.Tags {
		if t == tag {
			return
		}
	}
	m.Tags = append(m.Tags, tag)
}

func (m *Memory) ShareWith(agentID string) {
	for _, a := range m.SharedWith {
		if a == agentID {
			return
		}
	}
	m.SharedWith = append(m.SharedWith, agentID)
}

func (m *Memory) IsShared() bool {
	return len(m.SharedWith) > 0
}

func (m *Memory) IsSharedWith(agentID string) bool {
	for _, a := range m.SharedWith {
		if a == agentID {
			return true
		}
	}
	return false
}

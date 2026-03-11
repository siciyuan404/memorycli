package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/memorycli/memorycli/internal/config"
	"github.com/memorycli/memorycli/internal/memory"
	"github.com/spf13/cobra"
)

var (
	storeType string
	storeImportance float64
	storeTags string
	storeTTL int
	storeShare bool
)

var storeCmd = &cobra.Command{
	Use:   "store <content>",
	Short: "存储记忆",
	Long:  "存储新的记忆到 Redis",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content := args[0]

		agentID := config.GetAgentID()

		memType := memory.MemoryTypeKnowledge
		if storeType != "" {
			memType = memory.MemoryType(storeType)
		}

		m := memory.NewMemory(content, memType, agentID)
		m.SetImportance(storeImportance)

		if storeTags != "" {
			tags := strings.Split(storeTags, ",")
			for _, tag := range tags {
				m.AddTag(strings.TrimSpace(tag))
			}
		}

		if storeTTL > 0 {
			m.TTL = storeTTL
		}

		if storeShare {
			m.SharedWith = []string{"*"}
		}

		if err := store.Store(context.Background(), m); err != nil {
			return fmt.Errorf("failed to store memory: %w", err)
		}

		if jsonOutput {
			data, err := json.MarshalIndent(m, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal memory: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("记忆已存储: %s\n", m.ID)
		fmt.Printf("  类型: %s\n", m.Type)
		fmt.Printf("  重要性: %.2f\n", m.Importance)
		if len(m.Tags) > 0 {
			fmt.Printf("  标签: %v\n", m.Tags)
		}
		if m.IsShared() {
			fmt.Println("  已共享给所有 Agent")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.Flags().StringVarP(&storeType, "type", "t", "knowledge", "记忆类型 (knowledge|preference|conversation|task)")
	storeCmd.Flags().Float64VarP(&storeImportance, "importance", "i", 0.5, "重要性评分 (0-1)")
	storeCmd.Flags().StringVarP(&storeTags, "tags", "", "", "标签 (逗号分隔)")
	storeCmd.Flags().IntVarP(&storeTTL, "ttl", "", 0, "过期时间(秒)，0 表示永不过期")
	storeCmd.Flags().BoolVarP(&storeShare, "share", "s", false, "共享给所有 Agent")
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

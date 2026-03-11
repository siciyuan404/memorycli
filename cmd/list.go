package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/memorycli/memorycli/internal/memory"
	"github.com/spf13/cobra"
)

var (
	listLimit int
	listType string
	listAgent string
	listShared bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出记忆",
	Long:  "列出存储的记忆",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := memory.SearchOptions{
			Limit:      listLimit,
			AgentID:    listAgent,
			SharedOnly: listShared,
		}

		if listType != "" {
			opts.MemoryType = memory.MemoryType(listType)
		}

		memories, err := store.List(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("failed to list memories: %w", err)
		}

		if jsonOutput {
			data, err := json.MarshalIndent(memories, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal memories: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		if len(memories) == 0 {
			fmt.Println("没有找到记忆")
			return nil
		}

		fmt.Printf("共 %d 条记忆:\n\n", len(memories))
		for i, m := range memories {
			fmt.Printf("[%d] %s\n", i+1, m.ID)
			fmt.Printf("    Agent: %s\n", m.AgentID)
			fmt.Printf("    类型: %s\n", m.Type)
			fmt.Printf("    内容: %s\n", truncate(m.Content, 50))
			fmt.Printf("    重要性: %.2f\n", m.Importance)
			if len(m.Tags) > 0 {
				fmt.Printf("    标签: %v\n", m.Tags)
			}
			if m.IsShared() {
				fmt.Printf("    共享给: %v\n", m.SharedWith)
			}
			fmt.Printf("    创建时间: %s\n", m.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		return nil
	},
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	listCmd.Flags().StringVarP(&listType, "type", "t", "", "记忆类型 (knowledge|preference|conversation|task)")
	listCmd.Flags().StringVarP(&listAgent, "agent", "a", "", "指定 Agent ID")
	listCmd.Flags().BoolVarP(&listShared, "shared", "s", false, "只显示共享记忆")
}

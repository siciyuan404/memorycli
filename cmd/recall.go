package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/memorycli/memorycli/internal/memory"
	"github.com/spf13/cobra"
)

var (
	recallLimit int
	recallType string
	recallAgent string
)

var recallCmd = &cobra.Command{
	Use:   "recall <query>",
	Short: "检索记忆",
	Long:  "检索与查询相关的记忆",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		opts := memory.SearchOptions{
			Limit:      recallLimit,
			AgentID:    recallAgent,
			SharedOnly: false,
		}

		if recallType != "" {
			opts.MemoryType = memory.MemoryType(recallType)
		}

		memories, err := store.Search(context.Background(), query, opts)
		if err != nil {
			return fmt.Errorf("failed to recall memories: %w", err)
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
			fmt.Println("没有找到相关记忆")
			return nil
		}

		for i, m := range memories {
			fmt.Printf("[%d] %s\n", i+1, m.ID)
			fmt.Printf("    类型: %s\n", m.Type)
			fmt.Printf("    内容: %s\n", m.Content)
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

func init() {
	rootCmd.AddCommand(recallCmd)
	recallCmd.Flags().IntVarP(&recallLimit, "limit", "l", 10, "返回数量")
	recallCmd.Flags().StringVarP(&recallType, "type", "t", "", "记忆类型 (knowledge|preference|conversation|task)")
	recallCmd.Flags().StringVarP(&recallAgent, "agent", "a", "", "指定 Agent ID")
}

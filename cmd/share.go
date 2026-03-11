package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	shareTo string
)

var shareCmd = &cobra.Command{
	Use:   "share <memory_id>",
	Short: "共享记忆",
	Long:  "将记忆共享给其他 Agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		memoryID := args[0]

		if shareTo == "" {
			return fmt.Errorf("--to 参数是必填的")
		}

		targetAgents := strings.Split(shareTo, ",")
		for i, agent := range targetAgents {
			targetAgents[i] = strings.TrimSpace(agent)
		}

		if err := store.Share(context.Background(), memoryID, targetAgents); err != nil {
			return fmt.Errorf("failed to share memory: %w", err)
		}

		m, err := store.Get(context.Background(), memoryID)
		if err != nil {
			return fmt.Errorf("failed to get memory: %w", err)
		}

		if jsonOutput {
			data, err := json.MarshalIndent(m, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal memory: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("记忆 %s 已共享给: %v\n", memoryID, targetAgents)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringVarP(&shareTo, "to", "t", "", "目标 Agent ID (逗号分隔，必填)")
	shareCmd.MarkFlagRequired("to")
}

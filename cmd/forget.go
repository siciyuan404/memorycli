package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var forgetCmd = &cobra.Command{
	Use:   "forget <memory_id>",
	Short: "删除记忆",
	Long:  "删除指定的记忆",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		memoryID := args[0]

		m, err := store.Get(context.Background(), memoryID)
		if err != nil {
			return fmt.Errorf("failed to get memory: %w", err)
		}

		if err := store.Delete(context.Background(), memoryID); err != nil {
			return fmt.Errorf("failed to forget memory: %w", err)
		}

		fmt.Printf("记忆已删除: %s\n", memoryID)
		fmt.Printf("  内容: %s\n", truncate(m.Content, 50))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(forgetCmd)
}

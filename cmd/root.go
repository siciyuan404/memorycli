package cmd

import (
	"fmt"
	"os"

	"github.com/memorycli/memorycli/internal/config"
	"github.com/memorycli/memorycli/internal/memory"
	"github.com/memorycli/memorycli/internal/redis"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	jsonOutput bool
	redisURL string
	agentID string

	store *memory.Store
	redisClient *redis.Client
)

var rootCmd = &cobra.Command{
	Use:   "memorycli",
	Short: "MemoryCLI - AI Agent 记忆共享工具",
	Long: `MemoryCLI 是一个用 Golang 构建的 CLI 工具，
通过 skill 机制让多个 AI 编程助手实现记忆获取和共享。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.SetCfgFile(cfgFile)
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("failed to init config: %w", err)
		}

		cfg := config.GetConfig()
		if redisURL != "" {
			cfg.Redis.URL = redisURL
		}
		if agentID != "" {
			cfg.Agent.ID = agentID
		}

		var err error
		redisClient, err = redis.NewClient(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			return fmt.Errorf("failed to connect to redis: %w", err)
		}

		store = memory.NewStore(redisClient)

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if redisClient != nil {
			redisClient.Close()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径")
	rootCmd.PersistentFlags().StringVar(&redisURL, "redis", "", "Redis URL (默认 redis://localhost:6379)")
	rootCmd.PersistentFlags().StringVar(&agentID, "agent", "", "Agent ID")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "JSON 格式输出")
}

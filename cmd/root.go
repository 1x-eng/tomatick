package cmd

import (
	"fmt"
	"os"

	"github.com/1x-eng/tomatick/config"
	"github.com/1x-eng/tomatick/pkg/pomodoro"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "tomatick",
		Short: "A CLI Pomodoro timer with mem.ai integration",
		Run: func(cmd *cobra.Command, args []string) {
			pomo := pomodoro.NewTomatickMemento(cfg)
			pomo.StartCycle()
		},
	}
}

func Execute(cfg *config.Config) {
	rootCmd := NewRootCmd(cfg)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

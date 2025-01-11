package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bak",
	Short: "bak is a CLI tool for secure cloud storage integration",
	Long: `bak allows you to interact with cloud storage (S3, etc.) as if it were a local filesystem,
with built-in encryption and secure file handling capabilities.`,
}

func Execute() error {
	return rootCmd.Execute()
}

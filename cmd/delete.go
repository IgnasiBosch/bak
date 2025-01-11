package cmd

import (
	"fmt"
	"os"

	"bak/config"
	"bak/storage"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")
}

var deleteCmd = &cobra.Command{
	Use:   "delete [remote_path]",
	Short: "Delete a file from storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remotePath := args[0]
		force, _ := cmd.Flags().GetBool("force")

		cfg, err := config.LoadConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}

		if !force {
			fmt.Printf("Are you sure you want to delete %s? [y/N]: ", remotePath)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				return // not much we can do here
			}
			if response != "y" && response != "Y" {
				fmt.Println("Operation cancelled")
				return
			}
		}

		if err := deleteFile(remotePath, cfg); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted %s\n", remotePath)
	},
}

func deleteFile(remotePath string, cfg *config.Config) error {
	s3Client, err := storage.NewS3Client(
		cfg.Bucket,
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
	)
	if err != nil {
		return err
	}

	return s3Client.DeleteFile(remotePath)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"bak/config"
	"bak/crypto"
	"bak/storage"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download [remote_path] [local_path]",
	Short: "Download and decrypt a file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		remotePath := args[0]
		localPath := args[1]

		cfg, err := config.LoadConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}

		if err := downloadFile(remotePath, localPath, cfg); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}

		fmt.Printf("Successfully downloaded %s to %s\n", remotePath, localPath)
	},
}

func downloadFile(remotePath, localPath string, cfg *config.Config) error {
	s3Client, err := storage.NewS3Client(
		cfg.Bucket,
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
	)
	if err != nil {
		return err
	}

	// Download data
	data, err := s3Client.DownloadFile(remotePath, os.Stdout)
	if err != nil {
		return err
	}
	// Check if the file is encrypted
	if crypto.IsEncrypted(data) {
		fmt.Printf("\nDecrypting %s...\n", remotePath)

		// Generate decryption key
		key := []byte(cfg.SecretKey)
		if len(key) > 32 {
			key = key[:32]
		} else {
			paddedKey := make([]byte, 32)
			copy(paddedKey, key)
			key = paddedKey
		}

		// Decrypt the data
		data, err = crypto.Decrypt(data, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt: %v", err)
		}
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write to file
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil

}

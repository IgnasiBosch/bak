package cmd

import (
	"fmt"
	"os"

	"bak/config"
	"bak/crypto"
	"bak/storage"

	"github.com/spf13/cobra"
)

var encrypt bool

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&encrypt, "encrypt", "e", false, "Encrypt the file before uploading")
}

var uploadCmd = &cobra.Command{
	Use:   "upload [local_path] [remote_path]",
	Short: "Upload and encrypt a file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		localPath := args[0]
		remotePath := args[1]

		cfg, err := config.LoadConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}

		if err := uploadFile(localPath, remotePath, cfg); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}

		fmt.Printf("Successfully uploaded %s to %s\n", localPath, remotePath)
	},
}

func uploadFile(localPath, remotePath string, cfg *config.Config) error {
	// Read the file
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if encrypt {
		fmt.Printf("Encrypting %s...\n", localPath)

		// Generate encryption key from secret key
		key := []byte(cfg.SecretKey)
		if len(key) > 32 {
			key = key[:32]
		} else {
			paddedKey := make([]byte, 32)
			copy(paddedKey, key)
			key = paddedKey
		}

		// Encrypt the data
		data, err = crypto.Encrypt(data, key)
		if err != nil {
			return fmt.Errorf("failed to encrypt: %v", err)
		}
	}

	// Upload to S3
	s3Client, err := storage.NewS3Client(
		cfg.Bucket,
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
	)
	if err != nil {
		return err
	}

	if err := s3Client.UploadFile(localPath, remotePath, data, os.Stdout); err != nil {
		return fmt.Errorf("failed to upload: %v", err)
	}

	return nil
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if err != nil {
				return // Can't do anything if we can't write to stderr
			}
			os.Exit(1)
		}

		configDir := filepath.Join(homeDir, ".bak")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			if err != nil {
				return // Can't do anything if we can't write to stderr
			}
			os.Exit(1)
		}

		configFile := filepath.Join(configDir, "config.json")

		fmt.Print("Enter S3 endpoint URL: ")
		var endpoint string
		_, err = fmt.Scanln(&endpoint)
		if err != nil {
			return
		}

		fmt.Print("Enter access key: ")
		var accessKey string
		_, err = fmt.Scanln(&accessKey)
		if err != nil {
			return
		}

		fmt.Print("Enter secret key: ")
		var secretKey string
		_, err = fmt.Scanln(&secretKey)
		if err != nil {
			return
		}

		fmt.Print("Enter default bucket: ")
		var bucket string
		_, err = fmt.Scanln(&bucket)
		if err != nil {
			return
		}

		config := struct {
			Endpoint  string `json:"endpoint"`
			AccessKey string `json:"access_key"`
			SecretKey string `json:"secret_key"`
			Bucket    string `json:"bucket"`
		}{
			Endpoint:  endpoint,
			AccessKey: accessKey,
			SecretKey: secretKey,
			Bucket:    bucket,
		}

		data, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error encoding config: %v\n", err)
			if err != nil {
				return // Can't do anything if we can't write to stderr
			}
			os.Exit(1)
		}

		if err := os.WriteFile(configFile, data, 0600); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
			if err != nil {
				return // Can't do anything if we can't write to stderr
			}
			os.Exit(1)
		}

		fmt.Println("Configuration saved successfully!")
	},
}

package cmd

import (
	"fmt"
	"os"
	"time"

	"bak/config"
	"bak/storage"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolP("long", "l", false, "Show detailed information")
}

var lsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "List files and directories",
	Long:  `List files and directories in the specified path`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}

		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		showLong, _ := cmd.Flags().GetBool("long")

		if err := listS3Files(path, showLong, config); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if err != nil {
				return // not much we can do here
			}
			os.Exit(1)
		}
	},
}

func listS3Files(path string, longFormat bool, cfg *config.Config) error {
	if cfg.Bucket == "" {
		return fmt.Errorf("no default bucket configured. Run 'bak config set' first")
	}

	s3Client, err := storage.NewS3Client(
		cfg.Bucket,
		cfg.Endpoint,
		cfg.AccessKey,
		cfg.SecretKey,
	)
	if err != nil {
		return err
	}

	files, err := s3Client.ListFiles(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if longFormat {
			fmt.Printf("%8d %s %s\n",
				file.Size,
				file.LastModified.Format(time.RFC822),
				file.Name)
		} else {
			fmt.Println(file.Name)
		}
	}

	return nil
}

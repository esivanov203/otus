package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var ConfigFile string

func main() {
	rootCmd := &cobra.Command{
		Use:   "Scheduler",
		Short: "scheduler service entrypoint",
		RunE:  runScheduler,
	}

	rootCmd.PersistentFlags().StringVar(
		&ConfigFile,
		"config",
		"/etc/calendar/sender_config.yaml",
		"path to configuration file",
	)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show sender service version",
		Run:   printVersion,
	})

	if err := godotenv.Load(); err != nil {
		fmt.Printf(".env file not found or failed to load: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("initial error on: %v\n", err)
		os.Exit(1)
	}
}

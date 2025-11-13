package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/spf13/cobra"
)

var configFile string

func main() {
	rootCmd := &cobra.Command{
		Use:   "calendar",
		Short: "Calendar service entrypoint",
		RunE:  runServer,
	}

	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"/etc/calendar/config.yaml",
		"path to configuration file",
	)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show calendar service version",
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

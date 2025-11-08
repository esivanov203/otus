package main

import (
	"fmt"
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
		"/etc/calendar/config.toml",
		"path to configuration file",
	)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show calendar service version",
		Run:   printVersion,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

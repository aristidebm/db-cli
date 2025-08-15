package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"example.com/db/internal/client"
	"example.com/db/internal/config"
	"example.com/db/internal/types"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "db-cli",
		Short: "A CLI tool for managing database connections and running queries",
		Long:  "A comprehensive CLI tool that supports multiple database types including PostgreSQL, MySQL, SQLite, Redis, and more.",
	}

	// Initialize config
	config := &config.Config{}
	if err := config.Load(); err != nil {
		log.Fatalf("Warning: Could not load config: %v", err)
	}

	// Add commands
	rootCmd.AddCommand(
		createAddCommand(config),
		createPingCommand(config),
		createConnectCommand(config),
		createRunCommand(config),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func createAddCommand(config *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "add <source> <url>",
		Short: "Add a new data source",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(2), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.AddSource(args[0], args[1])
		},
	}
}

func createPingCommand(config *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "ping <source>",
		Short: "Ping a datasource",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			c, err := getClient(config, name, "")
			if err != nil {
				return err
			}
			return c.Ping()
		},
	}
}

func createConnectCommand(config *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "connect <source>",
		Short: "Connect to a datasource",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			c, err := getClient(config, name, "")
			if err != nil {
				return err
			}
			return c.Connect()
		},
	}
}

func createRunCommand(config *config.Config) *cobra.Command {
	var formatStr string = ""

	runCmd := &cobra.Command{
		Use:   "run <source> <query>",
		Short: "Run a query on a datasource",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(2), cobra.OnlyValidArgs),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// we need enumerations please !
			validFormat := []string{
				"json",
				"csv",
				"html",
				"md",
				"markdown",
				"latex",
				"asciidoc",
			}
			if formatStr != "" && !slices.Contains(validFormat, formatStr) {
				return fmt.Errorf("invalid format: '%s', valid choices are json, csv, html, md, markdown, latex, asciidoc", formatStr)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			c, err := getClient(config, name, types.Format(formatStr))
			if err != nil {
				return err
			}
			argsAfterDash := []string{}
			if dashIndex := cmd.ArgsLenAtDash(); dashIndex != -1 {
				argsAfterDash = append(argsAfterDash, args[dashIndex:]...)
			}
			return c.RunQuery(args[1], argsAfterDash...)
		},
	}
	runCmd.Flags().StringVarP(&formatStr, "format", "f", "", "output format")
	return runCmd
}

func getClient(config *config.Config, name string, format types.Format) (*client.Client, error) {
	source, err := config.GetSource(name)
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(source.URL)
	if err != nil {
		return nil, err
	}
	c.Format = format
	c.SourceConfig = source
	driver, err := config.GetDriver(c.Driver)
	if err == nil {
		c.DriverConfig = driver
	}
	return c, nil
}

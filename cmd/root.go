package cmd

import (
	"fmt"
	"os"

	"agentic/commerce/internal/app"

	"github.com/spf13/cobra"
)

var (
	configPath string

	rootCMD = &cobra.Command{
		Use:                "application",
		Short:              "Running best application in the world",
		PersistentPreRun:   preRun,
		PersistentPostRunE: postRun,
	}
)

func init() {
	cobra.OnInitialize(initialize)

	rootCMD.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yml", "Path of config file (using the default config if not specified)")

	rootCMD.AddCommand(serveCMD)
}

func initialize() {
	fmt.Println(app.Banner())
}

func preRun(_ *cobra.Command, _ []string) {
	fmt.Println("Starting up application...")
}

func postRun(_ *cobra.Command, _ []string) error {
	fmt.Println("Shutting down application...")
	return nil
}

// Execute executes the root command.
func Execute() {
	err := rootCMD.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
	}
}

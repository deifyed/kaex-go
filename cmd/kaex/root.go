package main

import (
	"os"
	"path"

	"github.com/deifyed/kaex-go/pkg/api"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "kaex",
		Short: "Kubernetes Application.yaml EXpander",
		Long:  "A tool for converting a simpler application.yaml into Kubernetes resources",
	}
)

func Execute() error {
	kaex := api.Kaex{
		Err: rootCmd.ErrOrStderr(),
		Out: rootCmd.OutOrStdout(),
		In:  rootCmd.InOrStdin(),

		TemplatesDirURL: "https://raw.githubusercontent.com/oslokommune/kaex/master/templates",
	}

	configPath, err := os.UserConfigDir()
	if err == nil {
		kaex.ConfigPath = path.Join(configPath, "kaex")
	}

	rootCmd.AddCommand(buildInitializeCommand(kaex))
	rootCmd.AddCommand(buildExpandCommand(kaex))
	rootCmd.AddCommand(buildSplitCommand(kaex))

	return rootCmd.Execute()
}

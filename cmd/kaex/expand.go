package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/deifyed/kaex-go/pkg/api"
	"github.com/spf13/cobra"
)

type ExpandOptions struct {
	Save    bool
	PodOnly bool
}

func buildExpandCommand(kaex api.Kaex) *cobra.Command {
	options := &ExpandOptions{
		Save:    false,
		PodOnly: false,
	}

	cmd := &cobra.Command{
		Use:     "expand",
		Aliases: []string{"exp", "x"},
		Short:   "expands an application.yaml from stdin (default)",
		Long:    "expands an application.yaml from stdin",
		RunE: func(_ *cobra.Command, args []string) error {
			app, err := api.ParseApplication(kaex.In)
			if err != nil {
				return err
			}

			var kubernetesResourceBuffer bytes.Buffer

			err = api.Expand(&kubernetesResourceBuffer, app, options.PodOnly)
			if err != nil {
				return err
			}

			err = write(kaex, app, &kubernetesResourceBuffer, options.Save)

			return err
		},
	}

	flags := cmd.Flags()

	flags.BoolVarP(&options.Save, "save", "s", false, "save the expanded Kubernetes resources to files")
	flags.BoolVarP(&options.PodOnly, "pod-only", "p", false, "create a pod resource instead of a deployment")

	return cmd
}

func write(kaex api.Kaex, app api.Application, reader io.Reader, asFile bool) error {
	var (
		writer io.Writer
		err    error
	)

	switch asFile {
	case false:
		writer = kaex.Out
	case true:
		filename := fmt.Sprintf("%s.yaml", app.Name)

		writer, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", filename, err)
		}
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		return fmt.Errorf("error writing: %w", err)
	}

	return nil
}

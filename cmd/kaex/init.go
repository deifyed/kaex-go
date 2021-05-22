package main

import (
	"bytes"
	"fmt"

	"github.com/deifyed/kaex-go/pkg/api"
	"github.com/spf13/cobra"
)

var templateName string

const templateNameHelp = `Define what template to use. Kaex will first look in ~/.config/kaex/templates, then in
https://github.com/oslokommune/kaex/tree/master/templates
`

const initializeHelpLong = `Scaffolds a template.

By default, kaex init will scaffold an application.yaml. It will first look in your ~/.config/kaex/templates folder for
an application.yaml file, then it will poll the /templates folder of the Kaex git repository.

By putting files in ~/.config/kaex/templates, you can configure what Kaex will scaffold. This enables you to create
different scaffolds for different environments and scaffold them for example by running kaex init -t <provider>.
`

func buildInitializeCommand(kaex api.Kaex) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "initialize",
		Aliases: []string{"init", "i"},
		Short:   "scaffolds a template",
		Long:    initializeHelpLong,
		RunE: func(_ *cobra.Command, args []string) error {
			var buffer bytes.Buffer

			err := api.FetchTemplate(kaex, &buffer, templateName)
			if err != nil {
				return fmt.Errorf("error fetching template: %w", err)
			}

			_, err = buffer.WriteTo(kaex.Out)
			if err != nil {
				return fmt.Errorf("error writing to output stream: %w", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&templateName, "template", "t", "application.yaml", templateNameHelp)

	return cmd
}

package main

import (
	"bytes"
	"fmt"
	"github.com/deifyed/kaex-go/pkg/api"
	"github.com/deifyed/kaex-go/pkg/extraction"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
)

const splitLongHelpText = `Splits a yaml file containing multiple Kubernetes resources.

Usage example:
$ kaex split <<EOF
apiVersion: v1
kind: Ingress

---
apiVersion: v1
kind: Service
EOF

$ cat file-containing-multiple-resources.yaml | kaex split

Result:
$ ls
. .. ingress.yaml service.yaml
`

var outputDirectory string

func buildSplitCommand(kaex api.Kaex) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "split",
		Aliases: []string{"s"},
		Short:   "split a yaml file with multiple Kubernetes resources into one file for each resource",
		Long:    splitLongHelpText,
		RunE: func(_ *cobra.Command, args []string) error {
			data, err := ioutil.ReadAll(kaex.In)
			if err != nil {
				return fmt.Errorf("reading input: %w", err)
			}

			fs := afero.Afero{Fs: afero.NewOsFs()}

			parts := bytes.Split(data, []byte("---"))

			for _, part := range parts {
				err = handlePart(fs, outputDirectory, part)
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&outputDirectory,
		"output-directory",
		"o",
		"./",
		"choose what directory Kaex should place the resource files",
	)

	return cmd
}

func handlePart(fs afero.Afero, outputDirectory string, part []byte) (err error) {
	kind, err := extraction.ExtractKind(part)
	if err != nil {
		return fmt.Errorf("extracting kind: %w", err)
	}

	targetPath := path.Join(outputDirectory, fmt.Sprintf("%s.yaml", kind))

	f, err := fs.Create(targetPath)
	if err != nil {
		return fmt.Errorf("creating new file: %w", err)
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.Write(part)
	if err != nil {
		return fmt.Errorf("writing part: %w", err)
	}

	return nil
}

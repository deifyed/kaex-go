package api

import (
	"bytes"
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

type IngressConfig struct {
	Annotations map[string]string
}

type Application struct {
	Name      string
	Namespace string

	Image           string
	Version         string
	ImagePullSecret string

	Url  string
	Port int32

	Replicas int32

	Environment map[string]string
	Volumes     []map[string]string

	Ingress IngressConfig
}

func ParseApplication(reader io.Reader) (Application, error) {
	var (
		app Application
		buf bytes.Buffer
	)

	_, err := io.Copy(&buf, reader)
	if err != nil {
		return Application{}, fmt.Errorf("error reading stdin: %w", err)
	}

	err = yaml.Unmarshal(buf.Bytes(), &app)
	if err != nil {
		return Application{}, fmt.Errorf("error parsing yaml: %w", err)
	}

	return app, nil
}

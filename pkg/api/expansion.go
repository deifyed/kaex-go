package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"sigs.k8s.io/yaml"
)

var emptyMatcher, _ = regexp.Compile("^.*: (null|{})$")
var statusMatcher, _ = regexp.Compile("^\\s*?status*:$")

func Expand(w io.Writer, app Application, podOnly bool) error {
	if app.Port != 0 {
		service, err := CreateService(app)
		if err != nil {
			return err
		}

		err = WriteCleanResource(w, service)
		if err != nil {
			return err
		}
	}

	if len(app.Volumes) != 0 {
		for _, volume := range app.Volumes {
			for path, size := range volume {
				volume, err := CreatePersistentVolume(app, path, size)
				if err != nil {
					return err
				}

				err = WriteCleanResource(w, volume)
				if err != nil {
					return err
				}
			}
		}
	}

	if app.Url != "" {
		ingress, err := CreateIngress(app)
		if err != nil {
			return err
		}
		err = WriteCleanResource(w, ingress)
		if err != nil {
			return err
		}
	}

	if podOnly == false {
		deployment, err := CreateDeployment(app)
		if err != nil {
			return err
		}
		err = WriteCleanResource(w, deployment)
		if err != nil {
			return err
		}
	} else {
		pod, err := CreatePod(app)
		if err != nil {
			return err
		}
		err = WriteCleanResource(w, pod)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteResource(w io.Writer, resource interface{}) error {
	serializedResource, err := yaml.Marshal(resource)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n---\n", serializedResource)
	if err != nil {
		return err
	}

	return nil
}

func WriteCleanResource(w io.Writer, resource interface{}) error {
	var buf bytes.Buffer

	err := WriteResource(&buf, resource)
	if err != nil {
		return err
	}

	result, err := cleanResources(buf)
	if err != nil {
		return err
	}

	_, err = w.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func cleanResources(buf bytes.Buffer) ([]byte, error) {
	content, err := ioutil.ReadAll(&buf)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	for _, item := range strings.Split(string(content), "\n") {
		if !emptyMatcher.MatchString(item) && !statusMatcher.MatchString(item) {
			result.Write([]byte(item + "\n"))
		}
	}

	return result.Bytes(), nil
}

package api

import (
	"net/url"
	"strings"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDefaultIngress() v1.Ingress {
	return v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "name",
			Labels:      nil,
			Annotations: nil,
		},
		Spec: v1.IngressSpec{
			Rules: make([]v1.IngressRule, 1),
		},
	}
}

func CreateIngress(app Application) (v1.Ingress, error) {
	hostUrl, err := url.Parse(app.Url)
	if err != nil {
		return v1.Ingress{}, err
	}

	ingress := generateDefaultIngress()
	ingress.ObjectMeta.Namespace = app.Namespace

	ingress.ObjectMeta.Name = app.Name
	ingress.ObjectMeta.Annotations = app.Ingress.Annotations

	ingress.Spec.Rules[0] = v1.IngressRule{
		Host: hostUrl.Host,
		IngressRuleValue: v1.IngressRuleValue{
			HTTP: &v1.HTTPIngressRuleValue{
				Paths: []v1.HTTPIngressPath{{
					Path: "/",
					Backend: v1.IngressBackend{
						Service: &v1.IngressServiceBackend{
							Name: app.Name,
							Port: v1.ServiceBackendPort{
								Number: 80,
							},
						},
					},
					PathType: pathTypeAsPtr(v1.PathTypePrefix),
				}},
			},
		},
	}

	if hostUrl.Scheme == "https" {
		ingress.Spec.TLS = []v1.IngressTLS{
			{
				Hosts: []string{
					hostUrl.Host,
				},
				SecretName: strings.Join([]string{
					app.Name,
					"tls",
				}, "-"),
			},
		}
	}

	return ingress, nil
}

func pathTypeAsPtr(p v1.PathType) *v1.PathType {
	return &p
}

package api_test

import (
	"bytes"
	"testing"

	"github.com/deifyed/kaex-go/pkg/api"
	"github.com/sebdah/goldie/v2"
)

func TestExpand(t *testing.T) {
	testCases := []struct {
		name string

		withApp api.Application
	}{
		{
			name: "Should produce expected output when expanding a simple app",

			withApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummygres",
				Version: "0.1.0",
			},
		},
		{
			name: "Should produce expected output when expanding app with optionals",

			withApp: api.Application{
				Name:            "dummy-app",
				Namespace:       "dummyns",
				Image:           "dummygres",
				Version:         "8.0.1",
				ImagePullSecret: "so-secret",
				Replicas:        3,
				Environment: map[string]string{
					"DUMMY_VAR":  "avalue",
					"DUMMY_HOST": "somehost",
				},
			},
		},
		{
			name: "Should produce expected output when expanding app with service",

			withApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummyredis",
				Version: "8.2.1",
				Port:    3000,
			},
		},
		{
			name: "Should produce expected output when expanding app with service and ingress",

			withApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummyredis",
				Version: "8.2.1",
				Port:    3000,
				Url:     "http://dummy.io",
			},
		},
		{
			name: "Should produce expected output when expanding app with annotated ingress",

			withApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummyredis",
				Version: "8.2.1",
				Port:    3000,
				Url:     "http://dummy.io",
				Ingress: api.IngressConfig{
					Annotations: map[string]string{
						"cert-manager.io/cluster-issuer": "letsencrypt-production",
					},
				},
			},
		},
		{
			name: "Should produce expected output when expanding app with service and tls enabled ingress",

			withApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummyredis",
				Version: "8.2.1",
				Port:    3000,
				Url:     "https://dummy.io",
			},
		},
		{
			name: "Should produce expected output when expanding app with volumes",

			withApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummygres",
				Version: "0.1.0",
				Volumes: []map[string]string{
					{"/etc/config": "4Gi"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")

			err := api.Expand(buf, tc.withApp, false)
			if err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)
			g.Assert(t, t.Name(), buf.Bytes())
		})
	}
}

type ExampleResource struct {
	Annotations map[string]string `json:"annotations"`
}

func TestWriteResources(t *testing.T) {
	testCases := []struct {
		name string

		with ExampleResource
	}{
		{
			name: "Should work should match snapshot",

			with: ExampleResource{
				Annotations: map[string]string{
					"alb.ingress.kubernetes.io/actions.ssl-redirect": "{\"Type\": \"redirect\", \"RedirectConfig\": { \"Protocol\": \"HTTPS\", \"Port\": \"443\", \"StatusCode\": \"HTTP_301\"}}",
					"alb.ingress.kubernetes.io/certificate-arn":      "arn:which:isnt:an:arn",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")

			err := api.WriteResource(buf, tc.with)
			if err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)
			g.Assert(t, t.Name(), buf.Bytes())
		})
	}
}
func TestWriteCleanResources(t *testing.T) {
	testCases := []struct {
		name string

		with ExampleResource
	}{
		{
			name: "Should match snapshot when cleaned",

			with: ExampleResource{
				Annotations: map[string]string{
					"alb.ingress.kubernetes.io/actions.ssl-redirect": "{\"Type\": \"redirect\", \"RedirectConfig\": { \"Protocol\": \"HTTPS\", \"Port\": \"443\", \"StatusCode\": \"HTTP_301\"}}",
					"alb.ingress.kubernetes.io/certificate-arn":      "arn:which:isnt:an:arn",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")

			err := api.WriteResource(buf, tc.with)
			if err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)
			g.Assert(t, t.Name(), buf.Bytes())
		})
	}
}

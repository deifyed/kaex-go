package api_test

import (
	"github.com/deifyed/kaex-go/pkg/api"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestParseApplication(t *testing.T) {
	testCases := []struct {
		name string

		withInput io.Reader

		expectApp api.Application
	}{
		{
			name: "Should work",

			withInput: strings.NewReader("name: dummy-app\nimage: dummygres\nversion: 0.1.0\n"),

			expectApp: api.Application{
				Name:    "dummy-app",
				Image:   "dummygres",
				Version: "0.1.0",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			resultApp, err := api.ParseApplication(tc.withInput)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expectApp, resultApp)
		})
	}
}

package extraction

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKind(t *testing.T) {
	testCases := []struct {
		name string

		with       []byte
		expectKind string

		expectFail error
	}{
		{
			name: "Should extract correct kind with valid input",

			with: []byte("apiVersion: v1\nkind: Ingress\nmetadata:\n\tname: Test"),

			expectKind: "ingress",
			expectFail: nil,
		},
		{
			name: "Should fail with ErrorNoMatch on invalid input",

			with: []byte("apiVersion: v1\nmetadata:\n\tname: Test\n"),

			expectFail: ErrorNoMatch,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractKind(tc.with)

			if tc.expectFail != nil {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(tc.expectFail, err))
			} else {
				assert.Nil(t, err)

				assert.Equal(t, tc.expectKind, result)
			}
		})
	}
}

package core

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/spf13/afero"
)

func TestGettingValueFromFile(t *testing.T) {
	testCases := []struct {
		name string

		withEnv   map[string]string
		withFiles map[string][]byte

		expectConfig func(cfg Config) bool
	}{
		{
			name: "Should work getting value from file when using @",

			withEnv: map[string]string{
				"DSN": "@/vault/secret/dsn",
			},
			withFiles: map[string][]byte{
				"/vault/secret/dsn": []byte("postgres://user:pass@host:5432/test"),
			},
			expectConfig: func(cfg Config) bool {
				return cfg.DSN.String() == "postgres://user:pass@host:5432/test"
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			for key, value := range tc.withFiles {
				f, err := fs.Create(key)
				assert.Equal(t, err, nil)

				_, err = f.Write(value)
				assert.Equal(t, err, nil)
			}

			mockEnvGetter := func(key string) string {
				value, ok := tc.withEnv[key]
				if !ok {
					return ""
				}

				return value
			}

			cfg := LoadConfig(fs, mockEnvGetter)

			assert.Equal(t, true, tc.expectConfig(cfg))
		})
	}
}

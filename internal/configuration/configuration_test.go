package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		envValue   *string
		defaultVal string
		want       string
	}{
		{name: "env set", envValue: ptr("from-env"), defaultVal: "fallback", want: "from-env"},
		{name: "env unset", envValue: nil, defaultVal: "fallback", want: "fallback"},
		{name: "env empty", envValue: ptr(""), defaultVal: "fallback", want: "fallback"},
		{name: "empty default", envValue: nil, defaultVal: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue == nil {
				t.Setenv("NAGARE_TEST_KEY", "")
			} else {
				t.Setenv("NAGARE_TEST_KEY", *tt.envValue)
			}

			assert.Equal(t, tt.want, GetEnvOrDefault("NAGARE_TEST_KEY", tt.defaultVal))
		})
	}
}

func TestLoad_UsesEnvDataDir(t *testing.T) {
	t.Setenv("NAGARE_DATA_DIR", "/tmp/nagare-custom")

	config := Load()

	assert.Equal(t, "/tmp/nagare-custom", config.DataDir)
}

func TestLoad_FallsBackToDefault(t *testing.T) {
	t.Setenv("NAGARE_DATA_DIR", "")

	config := Load()

	assert.Equal(t, DefaultConfiguration.DataDir, config.DataDir)
	assert.NotEmpty(t, config.DataDir)
}

func TestDefaultConfiguration_HasDataDir(t *testing.T) {
	assert.NotEmpty(t, DefaultConfiguration.DataDir)
}

func ptr(s string) *string { return &s }

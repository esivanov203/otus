package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		_, err := NewConfig("./noconfig.yaml")
		require.Error(t, err)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("invalid YAML", func(t *testing.T) {
		f, err := os.CreateTemp("", "config*.yaml")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.WriteString("::: this is not yaml :::")
		require.NoError(t, err)
		f.Close()

		_, err = NewConfig(f.Name())
		require.Error(t, err)
		require.Contains(t, err.Error(), "decoding")
	})

	t.Run("valid YAML", func(t *testing.T) {
		f, err := os.CreateTemp("", "config*.yaml")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		yamlData := `
logger:
  level: "INFO"
server:
  host: "localhost"
  port: 8080
storage:
  type: "memory"
`
		_, err = f.WriteString(yamlData)
		require.NoError(t, err)
		f.Close()

		cfg, err := NewConfig(f.Name())
		require.NoError(t, err)
		require.Equal(t, "INFO", cfg.Logger.Level)
		require.Equal(t, "localhost", cfg.Server.Host)
		require.Equal(t, 8080, cfg.Server.Port)
		require.Equal(t, "memory", cfg.Storage.Type)
	})
}

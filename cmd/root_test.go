package cmd

import (
	"os"
	"path/filepath"
	"testing"

	cfg "github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistentPreRunNormalizesInvalidConfiguredOutput(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", cfg.AppName)
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, cfg.ConfigFileName+".yaml")
	err = os.WriteFile(configPath, []byte("output_format: table\n"), 0o644)
	require.NoError(t, err)

	originalConfigFile := configFile
	configFile = configPath
	t.Cleanup(func() {
		configFile = originalConfigFile
	})

	cmd, _, err := rootCmd.Find([]string{"config", "get"})
	require.NoError(t, err)
	require.NotNil(t, cmd)

	err = rootCmd.PersistentPreRunE(cmd, []string{})
	require.NoError(t, err)
	assert.Equal(t, "yaml", cfg.Cfg.OutputFormat)
}

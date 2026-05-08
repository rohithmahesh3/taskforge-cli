package config

import (
	"os"
	"path/filepath"
	"testing"

	internalconfig "github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSetWorkspace(t *testing.T) {
	setupConfigTest(t)

	err := runSet(nil, []string{"workspace", "engineering"})
	require.NoError(t, err)
	assert.Equal(t, "engineering", internalconfig.Cfg.DefaultWorkspace)
}

func TestRunSetRejectsOutputKey(t *testing.T) {
	setupConfigTest(t)
	err := runSet(nil, []string{"output", "yaml"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown config key")
}

func setupConfigTest(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", internalconfig.AppName)
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	internalconfig.SetConfigFile(filepath.Join(configDir, internalconfig.ConfigFileName+".yaml"))
	internalconfig.Cfg = internalconfig.Config{}

	err = internalconfig.InitConfig()
	require.NoError(t, err)
}

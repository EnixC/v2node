package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

var serverEnvNames = []string{
	"V2NODE_API_HOST",
	"V2NODE_NODE_ID",
	"V2NODE_API_KEY",
	"V2NODE_TIMEOUT",
	"V2NODE_RETRY_COUNT",
	"V2NODE_LOG_LEVEL",
}

func TestLoadConfigFallsBackToFile(t *testing.T) {
	clearServerEnv(t)
	configPath := filepath.Join(t.TempDir(), "config.json")
	configJSON := `{
  "Nodes": [
    {
      "ApiHost": "https://panel.example.com/",
      "NodeID": 9,
      "ApiKey": "file-secret",
      "Timeout": 15
    }
  ]
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	c, envMode, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if envMode {
		t.Fatal("loadConfig() envMode = true, want false")
	}
	if len(c.NodeConfigs) != 1 || c.NodeConfigs[0].NodeID != 9 {
		t.Fatalf("NodeConfigs = %#v", c.NodeConfigs)
	}
}

func TestLoadConfigPrefersCompleteEnvironment(t *testing.T) {
	clearServerEnv(t)
	t.Setenv("V2NODE_API_HOST", "https://panel.example.com/")
	t.Setenv("V2NODE_NODE_ID", "11")
	t.Setenv("V2NODE_API_KEY", "env-secret")

	c, envMode, err := loadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if !envMode {
		t.Fatal("loadConfig() envMode = false, want true")
	}
	if len(c.NodeConfigs) != 1 || c.NodeConfigs[0].NodeID != 11 {
		t.Fatalf("NodeConfigs = %#v", c.NodeConfigs)
	}
}

func clearServerEnv(t *testing.T) {
	t.Helper()
	for _, name := range serverEnvNames {
		value, exists := os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			t.Fatalf("Unsetenv(%s): %v", name, err)
		}
		t.Cleanup(func() {
			if exists {
				_ = os.Setenv(name, value)
				return
			}
			_ = os.Unsetenv(name)
		})
	}
}

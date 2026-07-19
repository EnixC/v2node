package conf

import (
	"os"
	"strings"
	"testing"
)

var dockerEnvNames = []string{
	envAPIHost,
	envNodeID,
	envAPIKey,
	envTimeout,
	envRetryCount,
	envLogLevel,
}

func TestLoadFromEnvNotConfigured(t *testing.T) {
	clearDockerEnv(t)
	t.Setenv(envLogLevel, "warning")
	t.Setenv(envTimeout, "15")
	t.Setenv(envRetryCount, "1")

	c := New()
	configured, err := c.LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}
	if configured {
		t.Fatal("LoadFromEnv() configured = true, want false")
	}
}

func TestLoadFromEnvUsesDefaults(t *testing.T) {
	clearDockerEnv(t)
	t.Setenv(envAPIHost, " https://panel.example.com/api/ ")
	t.Setenv(envNodeID, "12")
	t.Setenv(envAPIKey, "secret")

	c := New()
	configured, err := c.LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}
	if !configured {
		t.Fatal("LoadFromEnv() configured = false, want true")
	}
	if len(c.NodeConfigs) != 1 {
		t.Fatalf("NodeConfigs length = %d, want 1", len(c.NodeConfigs))
	}

	node := c.NodeConfigs[0]
	if node.APIHost != "https://panel.example.com/api/" {
		t.Errorf("APIHost = %q", node.APIHost)
	}
	if node.NodeID != 12 {
		t.Errorf("NodeID = %d, want 12", node.NodeID)
	}
	if node.Key != "secret" {
		t.Errorf("Key = %q, want secret", node.Key)
	}
	if node.Timeout != DefaultNodeTimeout {
		t.Errorf("Timeout = %d, want %d", node.Timeout, DefaultNodeTimeout)
	}
	if node.RetryCount == nil || *node.RetryCount != DefaultNodeRetryCount {
		t.Errorf("RetryCount = %v, want %d", node.RetryCount, DefaultNodeRetryCount)
	}
}

func TestLoadFromEnvUsesOptionalValues(t *testing.T) {
	clearDockerEnv(t)
	t.Setenv(envAPIHost, "http://panel.example.com")
	t.Setenv(envNodeID, "7")
	t.Setenv(envAPIKey, "secret")
	t.Setenv(envTimeout, "30")
	t.Setenv(envRetryCount, "3")
	t.Setenv(envLogLevel, "WARN")

	c := New()
	configured, err := c.LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}
	if !configured {
		t.Fatal("LoadFromEnv() configured = false, want true")
	}

	node := c.NodeConfigs[0]
	if node.Timeout != 30 {
		t.Errorf("Timeout = %d, want 30", node.Timeout)
	}
	if node.RetryCount == nil || *node.RetryCount != 3 {
		t.Errorf("RetryCount = %v, want 3", node.RetryCount)
	}
	if c.LogConfig.Level != "warn" {
		t.Errorf("Log level = %q, want warn", c.LogConfig.Level)
	}
}

func TestLoadFromEnvRejectsPartialConfiguration(t *testing.T) {
	clearDockerEnv(t)
	t.Setenv(envAPIHost, "https://panel.example.com")

	c := New()
	configured, err := c.LoadFromEnv()
	if !configured {
		t.Fatal("LoadFromEnv() configured = false, want true")
	}
	if err == nil || !strings.Contains(err.Error(), "must be set together") {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}
}

func TestLoadFromEnvRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		value string
	}{
		{name: "API host", env: envAPIHost, value: "ftp://panel.example.com"},
		{name: "node ID", env: envNodeID, value: "0"},
		{name: "API key", env: envAPIKey, value: "   "},
		{name: "timeout", env: envTimeout, value: "0"},
		{name: "retry count", env: envRetryCount, value: "-1"},
		{name: "log level", env: envLogLevel, value: "trace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearDockerEnv(t)
			t.Setenv(envAPIHost, "https://panel.example.com")
			t.Setenv(envNodeID, "1")
			t.Setenv(envAPIKey, "secret")
			t.Setenv(tt.env, tt.value)

			c := New()
			configured, err := c.LoadFromEnv()
			if !configured {
				t.Fatal("LoadFromEnv() configured = false, want true")
			}
			if err == nil {
				t.Fatal("LoadFromEnv() error = nil, want validation error")
			}
		})
	}
}

func clearDockerEnv(t *testing.T) {
	t.Helper()
	for _, name := range dockerEnvNames {
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

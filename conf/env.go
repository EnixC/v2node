package conf

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	envAPIHost    = "V2NODE_API_HOST"
	envNodeID     = "V2NODE_NODE_ID"
	envAPIKey     = "V2NODE_API_KEY"
	envTimeout    = "V2NODE_TIMEOUT"
	envRetryCount = "V2NODE_RETRY_COUNT"
	envLogLevel   = "V2NODE_LOG_LEVEL"
)

// LoadFromEnv loads a single-node configuration when Docker-style environment
// variables are present. The bool result reports whether environment mode was
// selected.
func (p *Conf) LoadFromEnv() (bool, error) {
	apiHost, hasAPIHost := os.LookupEnv(envAPIHost)
	nodeIDValue, hasNodeID := os.LookupEnv(envNodeID)
	apiKey, hasAPIKey := os.LookupEnv(envAPIKey)

	if !hasAPIHost && !hasNodeID && !hasAPIKey {
		return false, nil
	}
	if !hasAPIHost || !hasNodeID || !hasAPIKey {
		return true, fmt.Errorf("%s, %s and %s must be set together", envAPIHost, envNodeID, envAPIKey)
	}

	apiHost = strings.TrimSpace(apiHost)
	if err := validateAPIHost(apiHost); err != nil {
		return true, err
	}
	if strings.TrimSpace(apiKey) == "" {
		return true, fmt.Errorf("%s must not be empty", envAPIKey)
	}

	nodeID, err := parseEnvInt(envNodeID, nodeIDValue, 1)
	if err != nil {
		return true, err
	}
	timeout, err := optionalEnvInt(envTimeout, DefaultNodeTimeout, 1)
	if err != nil {
		return true, err
	}
	retryCount, err := optionalEnvInt(envRetryCount, DefaultNodeRetryCount, 0)
	if err != nil {
		return true, err
	}

	if level, ok := os.LookupEnv(envLogLevel); ok && strings.TrimSpace(level) != "" {
		level = strings.ToLower(strings.TrimSpace(level))
		switch level {
		case "debug", "info", "warn", "warning", "error":
			p.LogConfig.Level = level
		default:
			return true, fmt.Errorf("%s has unsupported value %q", envLogLevel, level)
		}
	}

	p.NodeConfigs = []NodeConfig{
		{
			APIHost:    apiHost,
			NodeID:     nodeID,
			Key:        apiKey,
			Timeout:    timeout,
			RetryCount: intPtr(retryCount),
		},
	}
	return true, nil
}

func validateAPIHost(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("%s must be a valid HTTP or HTTPS URL", envAPIHost)
	}
	return nil
}

func optionalEnvInt(name string, fallback, minimum int) (int, error) {
	value, ok := os.LookupEnv(name)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	return parseEnvInt(name, value, minimum)
}

func parseEnvInt(name, value string, minimum int) (int, error) {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < minimum {
		return 0, fmt.Errorf("%s must be an integer greater than or equal to %d", name, minimum)
	}
	return parsed, nil
}

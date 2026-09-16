package config
import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerURL string `json:"serverUrl"`
	IdentityPath string `json:"identityPath"`
	LogLevel string `json:"logLevel"`
	HeartbeatInterval int `json:"heartbeatInterval"`
}

func defaults() Config {
	return Config {
		ServerURL: "https://localhost:8443",
		IdentityPath: "/etc/argus/agent/identity.json",
		LogLevel: "info",
		HeartbeatInterval: 30,
	}
}

func Load(path string) (Config, error) {
	cfg := defaults()

	if path != "" {
		if err := loadFile(path, &cfg); err != nil {
			return Config{}, fmt.Errorf("config: load file: %w", err)
		}
	}

	applyEnv(&cfg)

	return cfg, nil
}

func loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read config file: %w", err) 
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("ARGUS_SERVER_URL"); v != "" {
		cfg.ServerURL = v
	}
	if v := os.Getenv("ARGUS_IDENTITY_PATH"); v != "" {
		cfg.IdentityPath = v
	}
	if v := os.Getenv("ARGUS_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("ARGUS_HEARTBEAT_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.HeartbeatInterval = n
		}
	}
}
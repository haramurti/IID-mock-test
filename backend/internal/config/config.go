// Package config loads runtime settings from environment variables.
package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string // HTTP port, default 8080
	DataSource string // "mock" (default) or "expa"
	ExpaToken  string // only used when DataSource == "expa"; never sent to the browser
}

func Load() (Config, error) {
	cfg := Config{
		Port:       getenv("PORT", "8080"),
		DataSource: getenv("DATA_SOURCE", "mock"),
		ExpaToken:  os.Getenv("EXPA_TOKEN"),
	}
	switch cfg.DataSource {
	case "mock":
	case "expa":
		if cfg.ExpaToken == "" {
			return cfg, fmt.Errorf("DATA_SOURCE=expa requires EXPA_TOKEN")
		}
	default:
		return cfg, fmt.Errorf("invalid DATA_SOURCE %q (want mock or expa)", cfg.DataSource)
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	WyzieKey string `json:"wyzie_key"`
	TMDBKey  string `json:"tmdb_key"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".wyziesubs", "config.json")
}

func loadConfig() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// ensureConfig prompts for any missing API keys and saves them.
func ensureConfig(cfg *Config) error {
	needSave := false

	if cfg.WyzieKey == "" {
		fmt.Println(boldText("  Wyzie API Key Setup"))
		fmt.Println(grayText("  Get your free key at: ") + cyanText("https://sub.wyzie.io/redeem"))
		fmt.Println()
		key, err := prompt("  Enter Wyzie API key", "")
		if err != nil {
			return err
		}
		cfg.WyzieKey = key
		needSave = true
	}

	if cfg.TMDBKey == "" {
		fmt.Println()
		fmt.Println(boldText("  TMDB API Key Setup"))
		fmt.Println(grayText("  Get your free key at: ") + cyanText("https://www.themoviedb.org/settings/api"))
		fmt.Println()
		key, err := prompt("  Enter TMDB API key", "")
		if err != nil {
			return err
		}
		cfg.TMDBKey = key
		needSave = true
	}

	if needSave {
		if err := saveConfig(cfg); err != nil {
			printWarn("Could not save config: " + err.Error())
		} else {
			fmt.Println()
			printSuccess("Config saved to " + configPath())
			fmt.Println()
		}
	}

	return nil
}

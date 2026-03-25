// Package config provides configuration management for the mtg CLI.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration.
type Config struct {
	Projects      map[string]string            `json:"projects"`
	MailTemplates map[string]map[string]string `json:"mail_templates"`
}

// GetDefaultPath returns the default configuration file path.
func GetDefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(home, ".config", "mtg", "config.json")
}

// Load reads and parses the configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %w", err)
	}

	return &config, nil
}

// ResolvePrefix resolves a prefix from project name or returns the direct prefix.
func ResolvePrefix(project, prefix, configPath string) (string, error) {
	if project != "" {
		config, err := Load(configPath)
		if err != nil {
			return "", fmt.Errorf("設定ファイル読み込みエラー: %w", err)
		}

		p, ok := config.Projects[project]
		if !ok {
			fmt.Fprintf(os.Stderr, "エラー: プロジェクト '%s' が見つかりません\n", project)
			fmt.Fprintln(os.Stderr, "\n利用可能なプロジェクト:")
			for proj := range config.Projects {
				fmt.Fprintf(os.Stderr, "  - %s\n", proj)
			}
			return "", fmt.Errorf("プロジェクトが見つかりません")
		}
		return p, nil
	} else if prefix != "" {
		return prefix, nil
	}

	return "", fmt.Errorf("-project または -prefix フラグが必要です")
}

// Save writes the configuration to a file.
func Save(path string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("設定ファイル書き込みエラー: %w", err)
	}

	return nil
}

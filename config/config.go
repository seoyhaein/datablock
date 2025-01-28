package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MaxWatchCount int    `json:"MaxWatchCount"`
	RootDir       string `json:"rootDir"` // lustre-client 마운트된 폴더로 사용할 예정.
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// 추가적으로 필수 항목 검증
	if config.MaxWatchCount <= 0 {
		return nil, fmt.Errorf("missing or invalid 'MaxWatchCount' in configuration")
	}
	if config.RootDir == "" {
		return nil, fmt.Errorf("missing 'rootDir' in configuration")
	}

	return &config, nil
}

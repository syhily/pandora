package config

import (
	"log"
	"os"
	"path/filepath"
)

const (
	configFileName = "pandora.yml"
)

// DefaultConfigFile returns the default configuration root directory
func DefaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to read user home directory %v", err)
	}
	return filepath.Join(home, ".config", "pandora", configFileName)
}

package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

const (
	configFileName = "pandora.yml"
)

var config *PandoraConfig

// DefaultConfigFile returns the default configuration root directory
func DefaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to read user home directory %v", err)
	}
	return filepath.Join(home, ".config", "pandora", configFileName)
}

// PandoraConfig represents the structure of the pandora configuration file
type PandoraConfig struct {
	Upyun struct {
		Bucket   string `yaml:"bucket"`
		Operator string `yaml:"operator"`
		Password string `yaml:"password"`
	} `yaml:"upyun"`
	Asset struct {
		Scheme string `yaml:"scheme"`
		Domain string `yaml:"domain"`
		Path   struct {
			Image string `yaml:"image"`
			Music string `yaml:"music"`
		} `yaml:"path"`
	} `yaml:"asset"`
}

// ReadConfig loads the yaml based configuration file and deserializes it
func ReadConfig(configFile string) {
	// Initialize pandora config
	stat, err := os.Stat(configFile)
	if err != nil || !stat.Mode().IsRegular() {
		log.Fatalf(`It sees like you haven't provided a valid config file: %s.`, configFile)
	}

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Failed to load the config file from: %s.\nError: %v", configFile, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := yaml.NewDecoder(reader)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Invalid config file format or location %s.\nError: %v", configFile, err)
	}
}

func GetConfig() *PandoraConfig {
	return config
}

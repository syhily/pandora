package config

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"go.yaml.in/yaml/v4"
)

const (
	ConfigFileName = "gifts.yml"
)

// S3Configuration represents S3 storage configuration
type S3Configuration struct {
	Region          string `yaml:"region" prompt:"label=the s3 region,optional=true"`
	Endpoint        string `yaml:"endpoint" prompt:"label=the s3 endpoint,optional=true"`
	Bucket          string `yaml:"bucket" prompt:"label=the s3 bucket,required=true"`
	AccessKey       string `yaml:"accessKey" prompt:"label=the s3 access key,required=true"`
	AccessSecretKey string `yaml:"accessSecretKey" prompt:"label=the s3 access secret key,required=true"`
	PublicDomain    string `yaml:"publicDomain" prompt:"label=the s3 public domain,required=true,validate=http_suffix"`
}

// Retrieve implements aws.CredentialsProvider interface
func (c *S3Configuration) Retrieve(context.Context) (aws.Credentials, error) {
	if c.AccessKey == "" || c.AccessSecretKey == "" {
		return aws.Credentials{}, fmt.Errorf("no accessKey or AccessSecretKey is provided")
	}

	return aws.Credentials{
		AccessKeyID:     c.AccessKey,
		SecretAccessKey: c.AccessSecretKey,
	}, nil
}

// PandoraConfig represents the main configuration structure
type PandoraConfig struct {
	// The root file for storing the images
	ProjectRoot string `yaml:"projectRoot" prompt:"label=the project root,default=.,optional=true"`
	BlogRoot    string `yaml:"blogRoot" prompt:"label=the blog root directory,required=true"`
	Convert     struct {
		DefaultQuality int    `yaml:"defaultQuality" prompt:"label=the convert quality,default=75"`
		DefaultFormat  string `yaml:"defaultFormat" prompt:"label=the convert format,default=jpg,validate=image_format"`
	} `yaml:"convert"`
	ImageS3 S3Configuration `yaml:"imageS3"`
	MusicS3 S3Configuration `yaml:"musicS3"`
}

// DefaultConfigRoot returns the default configuration root directory
func DefaultConfigRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to read user home directory %v", err)
	}
	return filepath.Join(home, ".config", "pandora")
}

// ReadConfig loads the yaml based configuration file and deserializes it
func ReadConfig(configPath string) *PandoraConfig {
	// Initialize pandora config
	stat, err := os.Stat(configPath)
	if err != nil || !stat.IsDir() {
		log.Fatalf(`It sees like you haven't config the tool.\nExecute the command "pandora config" for initializing.`)
	}

	file, err := os.Open(filepath.Join(configPath, ConfigFileName))
	if err != nil {
		log.Fatalf("Failed to load the config file from: %s.\nError: %v", configPath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := yaml.NewDecoder(reader)

	var c PandoraConfig
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatalf("Invalid config file format or location %s.\nError: %v", configPath, err)
	}
	return &c
}

// WriteConfig writes the configuration to a file
func WriteConfig(configPath string, config *PandoraConfig) error {
	stat, err := os.Stat(configPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(configPath, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("failed to create the config path %s: %w", configPath, err)
		}
	} else if err != nil || !stat.IsDir() {
		return fmt.Errorf("invalid config path %s", configPath)
	}

	configFile := filepath.Join(configPath, ConfigFileName)
	file, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("failed to create config file %s: %w", configFile, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("failed to generate the configuration file: %w", err)
	}

	return nil
}

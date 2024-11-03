package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HttpServer holds configuration for the HTTP server.
type HttpServer struct {
	Address string `yaml:"address" env-required:"true"` // Include a YAML tag if necessary for the address.
}

// Config represents the application configuration structure.
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"production"` // Environment (e.g., development, production).
	StoragePath string     `yaml:"storage_path" env-required:"true"`                           // Path to storage.
	HTTPServer  HttpServer `yaml:"http_server"`                                                // HTTP server configuration.
}

// MustLoad loads the configuration from a file or environment variables.
// It panics if loading fails, so it should only be used at startup.
func MustLoad() *Config {
	// Determine the configuration path from the environment variable or command-line flag.
	configPath := os.Getenv("CONFIG_PATH")

	// If CONFIG_PATH is not set, look for command-line flag.
	if configPath == "" {
		flags := flag.String("config-path", "", "path to config file")
		flag.Parse()
		configPath = *flags

		// Ensure that a config path is provided either via environment variable or command-line flag.
		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set or config path is not provided")
		}
	}

	// Check if the config file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	// Load configuration using cleanenv.
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

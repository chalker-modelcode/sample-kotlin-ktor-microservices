package config

import "github.com/caarlos0/env/v11"

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

// Load reads configuration from environment variables and returns it.
// It fails fast if any required variable is missing or cannot be parsed.
func Load() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}

package config

// Env returns the value of the environment variable named by the key.
type Env string

const (
	// EnvDevelopment is the development environment
	EnvDevelopment Env = "development"
	// EnvStaging is the staging environment
	EnvStaging Env = "staging"
	// EnvProduction is the production environment
	EnvProduction Env = "production"
)

// Config is the configuration struct
type Config struct {
	Port        *string `env:"PORT"`
	DatabaseURL string  `env:"DATABASE_URL"`
	Env         string  `env:"ENV"`
}

// GetEnv returns the current environment
func (c *Config) GetEnv() Env {
	return Env(c.Env)
}

// Instance is the global configuration
var Instance *Config

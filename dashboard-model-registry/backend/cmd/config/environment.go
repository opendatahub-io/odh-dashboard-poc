package config

import "errors"

// Environment (development|staging|production)
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

type Config struct {
	Port int
	Env  Environment
}

func (e *Environment) String() string {
	return string(*e)
}

func (e *Environment) Set(value string) error {
	switch value {
	case "development", "staging", "production":
		*e = Environment(value)
		return nil
	default:
		return errors.New("invalid environment: must be one of 'development', 'staging', or 'production'")
	}
}

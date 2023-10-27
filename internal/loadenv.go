package internal

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	env := os.Getenv("GO_ENV")
	envName := ".env.development"

	if env == "production" {
		envName = ".env.production"
	}

	err := godotenv.Load(envName)
	if err != nil {
		return err
	}

	return nil
}

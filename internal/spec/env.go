package spec

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from the .env files for tests run
func LoadEnv() error {
	rootDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "test"
	}

	envFiles := []string{
		filepath.Join(rootDir, ".env"),
		filepath.Join(rootDir, fmt.Sprintf(".env.%s", env)),
		filepath.Join(rootDir, fmt.Sprintf(".env.%s.local", env)),
	}
	for _, file := range envFiles {
		if err = godotenv.Overload(file); err != nil {
			log.Printf("%s file not found, skipping", file)
		}
	}

	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err = os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		if _, err = os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("could not find project root")
		}

		dir = parentDir
	}
}

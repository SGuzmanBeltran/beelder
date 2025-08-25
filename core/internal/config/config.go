package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatal("Could get ENV variable ", key)
	return ""
}

func LoadEnv(configPath string) error {
	// First try loading a .env located in the same directory as this source file.
	// This is handy during development when working directory may differ.
	if err := LoadEnvFromThisFile(configPath); err == nil {
		return nil
	}

	// Start from current directory and keep going up until we find .env or hit root
	// dir, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	// for {
	// 	if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
	// 		return godotenv.Load(filepath.Join(dir, ".env"))
	// 	}

	// 	parent := filepath.Dir(dir)
	// 	if parent == dir {
	// 		return errors.New(".env file not found")
	// 	}
	// 	dir = parent
	// }
	return nil
}

// LoadEnvFromThisFile loads a .env file located in the same directory
// as this source file (useful during development). It determines the
// file path using runtime.Caller and returns an error if the .env is
// not found or can't be loaded.
func LoadEnvFromThisFile(configPath string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("unable to determine caller file location")
	}

	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, configPath ,".env")

	if _, err := os.Stat(envPath); err != nil {
		return err
	}

	return godotenv.Load(envPath)
}

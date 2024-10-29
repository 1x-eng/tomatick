package config

import (
	"os"
	"path/filepath"
)

func getDefaultContextDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory cannot be determined
		return filepath.Join(".", ".tomatick", "context")
	}
	return filepath.Join(homeDir, ".tomatick", "context")
}

func ensureDirectoryExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

package pathutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveFilePath(inputPath string) (string, error) {
	if filepath.IsAbs(inputPath) {
		return inputPath, nil
	}

	if strings.HasPrefix(inputPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get home directory: %v", err)
		}
		return filepath.Join(homeDir, inputPath[1:]), nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %v", err)
	}

	return filepath.Join(pwd, inputPath), nil
}

func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

func IsRelativePath(path string) bool {
	return !filepath.IsAbs(path) && !strings.HasPrefix(path, "~")
}

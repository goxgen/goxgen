package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// ExecCommand executes command in the specified directory
func ExecCommand(dir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveFromDirByPatterns removes files from the specified directory by pattern
func RemoveFromDirByPatterns(patterns ...string) error {
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, f := range files {
			if err := os.RemoveAll(f); err != nil {
				return err
			}
		}
	}
	return nil
}

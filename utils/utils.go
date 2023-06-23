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

// RemoveFromDirByPatterns removes files and dirs from the specified directory by pattern
// eg ./*/dir1/*_qwe.go
//
//	./dir1/*_dir2
func RemoveFromDirByPatterns(patterns ...string) error {
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = os.RemoveAll(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

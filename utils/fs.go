package utils

import (
	"errors"
	"os"
)

func MkdirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}

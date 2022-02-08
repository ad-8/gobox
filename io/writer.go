package io

import (
	"os"
)

// SimpleWrite writes data to the named file, creating it if necessary.
func SimpleWrite(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0666); err != nil {
		return err
	}
	return nil
}

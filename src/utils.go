package src

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

func changeFileExtension(fileName, newExtension string) string {
	// Get the base name without extension
	baseName := filepath.Base(fileName)

	// Remove the old extension
	withoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]

	// Append the new extension
	return withoutExt + newExtension
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func mapToSlice[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func approximatelyEqual(a, b, epsilon float32) bool {
	return math.Abs(float64(a-b)) < float64(epsilon)
}

func removeAllFiles(dir string) error {
	// Open the directory
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Iterate through the directory entries
	for _, entry := range dirEntries {
		// Create the full path to the file or subdirectory
		fullPath := filepath.Join(dir, entry.Name())

		// If it's a directory, call the removeAllFiles function recursively
		if entry.IsDir() {
			err := removeAllFiles(fullPath)
			if err != nil {
				return err
			}
			// Remove the empty directory after its contents are deleted
			err = os.Remove(fullPath)
			if err != nil {
				return fmt.Errorf("failed to remove directory: %w", err)
			}
		} else {
			// If it's a file, remove it
			err := os.Remove(fullPath)
			if err != nil {
				return fmt.Errorf("failed to remove file: %w", err)
			}
		}
	}
	return nil
}

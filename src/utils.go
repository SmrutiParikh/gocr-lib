package src

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

func ChangeFileExtension(fileName, newExtension string) string {
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

func MapToSlice[K comparable, V any](m map[K]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func ApproximatelyEqual(a, b, epsilon float32) bool {
	return math.Abs(float64(a-b)) < float64(epsilon)
}

func RemoveAllFiles(dir string) error {
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
			err := RemoveAllFiles(fullPath)
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

// Function to read class names (from coco.names file)
func ReadClasses(classFile string) ([]string, error) {
	var classes []string
	// Load class labels from file
	file, err := os.Open(classFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		classes = append(classes, strings.TrimSpace(scanner.Text()))
	}
	return classes, scanner.Err()
}

func GetClassIndexAndConfidence(data gocv.Mat) (int, float32) {
	confidence := float32(0)
	index := 5
	for j := 5; j < data.Cols(); j++ {
		if data.GetFloatAt(0, j) > confidence {
			index = j
			confidence = data.GetFloatAt(0, j)
		}

	}
	return index - 5, confidence
}

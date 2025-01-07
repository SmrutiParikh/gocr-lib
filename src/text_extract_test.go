package src

import (
	"os"
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestExtractText(t *testing.T) {
	imageFiles := map[string]string{
		"input-image.png":  "eng",
		"crooked-scan.png": "eng",
		"bill.jpg":         "eng",
		"japanese.png":     "jpn"}

	expectedOutputs := map[string]string{
		"input-image.png":  "../samples/extracted-text/input-image.txt",
		"crooked-scan.png": "../samples/extracted-text/crooked-scan.txt",
		"bill.jpg":         "../samples/extracted-text/bill.txt",
		"japanese.png":     "../samples/extracted-text/japanese.txt",
	}

	algo := NewTextExtractionAlgorithm()

	// Loop over image files and verify the output
	for fileName, lang := range imageFiles {
		// Extract text from the image
		extractedText := algo.Execute("../samples/"+fileName, lang)

		// Get the expected output for the file
		expectedText, err := os.ReadFile(expectedOutputs[fileName])
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}

		// Compare the extracted text with the expected text
		if extractedText != string(expectedText) {
			t.Errorf("Test failed for file: %s. Expected: \n%s\n\n, but got: \n%s", fileName, expectedText, extractedText)
		}
	}
}

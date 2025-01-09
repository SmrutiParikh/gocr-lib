package doc

import (
	"go-ocr/src"
	"os"
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestTextExtraction(t *testing.T) {
	inputFiles := map[string]string{
		"input-image.png":  "eng",
		"crooked-scan.png": "eng",
		"bill.jpg":         "eng",
		"japanese.png":     "jpn"}

	pte := NewPlainTextExtractor()

	// Loop over image files and verify the output
	for fileName, lang := range inputFiles {
		// Extract text from the image
		extractedText := pte.Execute("../../samples/documents/"+fileName, lang)

		// Get the expected output for the file
		expectedText, err := os.ReadFile("../../output/test/extracted-text/" + src.ChangeFileExtension(fileName, ".txt"))
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}

		// Compare the extracted text with the expected text
		if extractedText != string(expectedText) {
			t.Errorf("Test failed for file: %s. Expected: \n%s\n\n, but got: \n%s", fileName, expectedText, extractedText)
		}
	}
}

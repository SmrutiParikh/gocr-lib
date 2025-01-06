package src

import (
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestExtractTextHOCR(t *testing.T) {
	removeAllFiles("../samples/generated-pdf")

	imageFiles := []string{"input-image.png", "crooked-scan.png", "bill.jpg"}

	algo := NewTextExtractionHOCRAlgorithm()

	// Loop over image files and verify the output
	for _, fileName := range imageFiles {
		algo.Execute("../samples/"+fileName, "../samples/generated-pdf/")
		if !FileExists("../samples/generated-pdf/" + changeFileExtension(fileName, ".pdf")) {
			t.Fatalf("Output pdf not generated for file %s", fileName)
		}
	}
}

package src

import (
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestExtractTextHOCR(t *testing.T) {
	removeAllFiles("../samples/generated-pdf")

	imageFiles := []string{"input-image.png", "crooked-scan.png", "bill.jpg"}

	algo := NewTextExtractionHOCRAlgorithm("../fonts/")

	// Loop over image files and verify the output
	for _, fileName := range imageFiles {
		outfilePath, err := algo.Execute("../samples/"+fileName, "../samples/generated-pdf/")
		if err != nil || !FileExists(*outfilePath) {
			t.Fatalf("Output pdf not generated for file %s", fileName)
		}
	}
}

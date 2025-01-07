package src

import (
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestExtractTextHOCR(t *testing.T) {
	removeAllFiles("../samples/generated-pdf")

	imageFiles := map[string]string{
		"input-image.png":  "eng",
		"crooked-scan.png": "eng",
		"bill.jpg":         "eng",
		"japanese.png":     "jpn"}

	algo := NewTextExtractionHOCRAlgorithm("../fonts/")

	// Loop over image files and verify the output
	for fileName, lang := range imageFiles {
		outfilePath, err := algo.Execute("../samples/"+fileName, lang, "../samples/generated-pdf/")
		if err != nil || !FileExists(*outfilePath) {
			t.Fatalf("Output pdf not generated for file %s", fileName)
		}
	}
}

package doc

import (
	"go-ocr/src"
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestHOCRTestExtraction(t *testing.T) {
	src.RemoveAllFiles("../../output/test/generated-pdf")

	inputFiles := map[string]string{
		"input-image.png":  "eng",
		"crooked-scan.png": "eng",
		"bill.jpg":         "eng",
		"japanese.png":     "jpn"}

	hte := NewHOCRTextExtractor("../../fonts/")

	// Loop over image files and verify the output
	for fileName, lang := range inputFiles {
		outfilePath, err := hte.Execute("../../samples/documents/"+fileName, lang,
			"../../output/test/generated-pdf/")
		if err != nil || !src.FileExists(*outfilePath) {
			t.Fatalf("Output pdf not generated for file %s", fileName)
		}
	}
}

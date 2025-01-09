package videos

import (
	"go-ocr/src"
	"testing"
)

// Unit test for checking text extraction from multiple images
func TestVideoObjectDetection(t *testing.T) {
	src.RemoveAllFiles("../../output/test/generated-videos")

	inputFiles := []string{"marathon.mp4"}

	vod := NewVideoObjectDetector(
		"../../models/yolov3.weights",
		"../../models/yolov3.cfg",
		"../../models/coco.names")

	// Loop over image files and verify the output
	for _, fileName := range inputFiles {
		outfilePath, err := vod.Execute("../../samples/videos/"+fileName,
			"../../output/test/generated-videos/", false)
		if err != nil || !src.FileExists(*outfilePath) {
			t.Fatalf("Output video not generated for file %s", fileName)
		}
	}
}

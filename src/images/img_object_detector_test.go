package images

import (
	"go-ocr/src"
	"testing"
)

// Test function
func TestImageObjectDetection(t *testing.T) {
	// Set up the image files and expected output
	inputFiles := []string{"kangaroo_horse.png", "traffic.jpg", "zebra_horse.jpg", "zebra.jpg"}
	expectedOutput := map[string][]DetectedObject{
		"kangaroo_horse.png": {
			{ClassName: "horse", Confidence: 0.95},
		},
		"traffic.jpg": {
			{ClassName: "car", Confidence: 0.80},
			{ClassName: "traffic light", Confidence: 0.83},
		},
		"zebra_horse.jpg": {
			{ClassName: "horse", Confidence: 0.67},
			{ClassName: "zebra", Confidence: 1.00},
		},
		"zebra.jpg": {
			{ClassName: "zebra", Confidence: 0.93},
		},
	}

	iod := NewImageObjectDetector(
		"../../models/yolov3.weights",
		"../../models/yolov3.cfg",
		"../../models/coco.names")

	defer iod.Close()

	// Iterate over each file and test detection
	for _, fileName := range inputFiles {
		detectedObjects := iod.Execute("../../samples/images/" + fileName)

		// Check if the output matches the expected values
		for i, detectedObject := range detectedObjects {
			expectedObj := expectedOutput[fileName][i]
			if detectedObject.ClassName != expectedObj.ClassName || !src.ApproximatelyEqual(detectedObject.Confidence, expectedObj.Confidence, 0.005) {
				t.Errorf("For file %s, expected %s with confidence %.2f, but got %s with confidence %.2f",
					fileName, expectedObj.ClassName, expectedObj.Confidence, detectedObject.ClassName, detectedObject.Confidence)
			}
		}
	}
}

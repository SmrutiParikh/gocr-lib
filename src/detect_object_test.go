package src

import (
	"testing"
)

// Test function
func TestObjectDetection(t *testing.T) {
	// Set up the image files and expected output
	imageFiles := []string{"kangaroo_horse.png", "traffic.jpg", "zebra_horse.jpg", "zebra.jpg"}
	expectedOutput := map[string][]DetectedObject{
		"kangaroo_horse.png": {
			{ClassName: "horse", Confidence: 0.96},
		},
		"traffic.jpg": {
			{ClassName: "car", Confidence: 0.81},
			{ClassName: "traffic light", Confidence: 0.83},
		},
		"zebra_horse.jpg": {
			{ClassName: "horse", Confidence: 0.94},
			{ClassName: "zebra", Confidence: 1.00},
		},
		"zebra.jpg": {
			{ClassName: "zebra", Confidence: 0.93},
		},
	}

	algo := NewObjectDetectionAlgorithm("../models/yolov3.weights", "../models/yolov3.cfg", "../models/coco.names")
	defer algo.Close()

	// Iterate over each file and test detection
	for _, fileName := range imageFiles {
		detectedObjects := algo.Execute("../samples/" + fileName)

		// Check if the output matches the expected values
		for i, detectedObject := range detectedObjects {
			expectedObj := expectedOutput[fileName][i]
			if detectedObject.ClassName != expectedObj.ClassName || !approximatelyEqual(detectedObject.Confidence, expectedObj.Confidence, 0.005) {
				t.Errorf("For file %s, expected %s with confidence %.2f, but got %s with confidence %.2f",
					fileName, expectedObj.ClassName, expectedObj.Confidence, detectedObject.ClassName, detectedObject.Confidence)
			}
		}
	}
}

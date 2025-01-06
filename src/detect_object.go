package src

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strings"

	"gocv.io/x/gocv"
)

type ObjectDetectionAlgo struct {
	net     *gocv.Net
	classes []string
}

func NewObjectDetectionAlgorithm(weightsFile string, configFile string, classFile string) *ObjectDetectionAlgo {
	// Load the pre-trained YOLOv3 model, configuration, and class labels
	net := gocv.ReadNet(weightsFile, configFile)
	if net.Empty() {
		fmt.Println("Error reading the network model")
		os.Exit(1)
	}

	// Load class labels (COCO class labels)
	classes, err := readClasses(classFile)
	if err != nil {
		fmt.Println("Error reading class names:", err)
		os.Exit(1)
	}

	return &ObjectDetectionAlgo{&net, classes}
}

func (algo *ObjectDetectionAlgo) Close() error {
	return algo.net.Close()
}

func (algo *ObjectDetectionAlgo) Execute(fileName string) []DetectedObject {
	// Load the input image
	img := gocv.IMRead(fileName, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading the image")
		return nil
	}
	defer img.Close()

	// Prepare the image for the network (convert to blob)
	blob := gocv.BlobFromImage(img, 0.00392, image.Pt(416, 416), gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()

	// Set the input to the network
	algo.net.SetInput(blob, "")

	// Get the output layer names
	layerNames := algo.net.GetLayerNames()

	// Run a forward pass: Get the output of the network
	// This returns a list of layers where the detections will come from
	detections := algo.net.ForwardLayers(layerNames)

	// Process the output to detect objects
	detectedObjects := make(map[int]DetectedObject, 0)

	// Process detections
	for _, detection := range detections {
		// Loop through the detection matrix (each object detected)
		// Each detection result contains a vector with the structure:
		// [centerX, centerY, width, height, confidence, classID, ...]
		for i := 0; i < detection.Rows(); i++ {
			// Extract the detection information
			confidence := detection.GetFloatAt(i, 4)
			if confidence > 0.5 { // Confidence threshold
				// Class ID (class of the detected object)
				classID := getClassIndex(detection, i)
				className := algo.classes[classID]
				if _, ok := detectedObjects[classID]; ok {
					detectedObjects[classID] = DetectedObject{className, confidence}
				} else if detectedObjects[classID].Confidence < confidence {
					detectedObjects[classID] = DetectedObject{className, confidence}
				}
			}
		}
	}

	return mapToSlice(detectedObjects)
}

type DetectedObject struct {
	ClassName  string
	Confidence float32
}

func getClassIndex(detection gocv.Mat, row int) int {
	max := float32(0)
	index := 5
	for j := 5; j < detection.Cols(); j++ {
		if detection.GetFloatAt(row, j) > max {
			index = j
			max = detection.GetFloatAt(row, j)
		}

	}
	return index - 5
}

// Function to read class names (from coco.names file)
func readClasses(classFile string) ([]string, error) {
	var classes []string
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

package images

import (
	"fmt"
	"image"
	"os"

	"go-ocr/src"

	"gocv.io/x/gocv"
)

type ImageObjectDetection struct {
	net        *gocv.Net
	classes    []string
	layerNames []string
}

func NewImageObjectDetector(weightsFile string, configFile string, classFile string) *ImageObjectDetection {
	// Load the pre-trained YOLOv3 model, configuration, and class labels
	net := gocv.ReadNet(weightsFile, configFile)
	if net.Empty() {
		fmt.Println("Error reading the network model")
		os.Exit(1)
	}

	// Load class labels (COCO class labels)
	classes, err := src.ReadClasses(classFile)
	if err != nil {
		fmt.Println("Error reading class names:", err)
		os.Exit(1)
	}

	// Get the output layer names
	layerNames := net.GetLayerNames()

	return &ImageObjectDetection{&net, classes, layerNames}
}

func (iod *ImageObjectDetection) Close() error {
	return iod.net.Close()
}

func (iod *ImageObjectDetection) Execute(fileName string) []DetectedObject {
	// Load the input image
	img := gocv.IMRead(fileName, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading the image")
		return nil
	}
	defer img.Close()

	// Prepare the image for the network (convert to blob)
	blob := gocv.BlobFromImage(img, 0.00392, image.Pt(416, 416),
		gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()

	// Set the input to the network
	iod.net.SetInput(blob, "")

	// Run a forward pass: Get the output of the network
	// This returns a list of layers where the layerOutputs will come from
	layerOutputs := iod.net.ForwardLayers(iod.layerNames)

	// Process the output to detect objects
	detectedObjects := make(map[int]DetectedObject, 0)

	// Process detections
	for _, output := range layerOutputs {
		// Loop through the detection matrix (each object detected)
		// Each detection result contains a vector with the structure:
		// [centerX, centerY, width, height, confidence, classID, ...]
		for i := 0; i < output.Rows(); i++ {
			data := output.Row(i)
			// Extract the detection information
			// Class ID (class of the detected object)
			classId, confidence := src.GetClassIndexAndConfidence(data)
			if confidence > 0.5 { // Confidence threshold
				className := iod.classes[classId]
				if _, ok := detectedObjects[classId]; ok {
					detectedObjects[classId] = DetectedObject{className, confidence}
				} else if detectedObjects[classId].Confidence < confidence {
					detectedObjects[classId] = DetectedObject{className, confidence}
				}
			}
		}
	}

	return src.MapToSlice(detectedObjects)
}

type DetectedObject struct {
	ClassName  string
	Confidence float32
}

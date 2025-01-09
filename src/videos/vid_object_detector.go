package videos

import (
	"fmt"
	"go-ocr/src"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gocv.io/x/gocv"
)

type VideoObjectDetection struct {
	net               *gocv.Net
	classes           []string
	layersNamesOutput []string
}

func NewVideoObjectDetector(weightsFile string, configFile string, classFile string) *VideoObjectDetection {
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

	// Perform forward pass
	layerNames := net.GetLayerNames()
	outputLayers := net.GetUnconnectedOutLayers()
	var layersNamesOutput []string
	for _, id := range outputLayers {
		layersNamesOutput = append(layersNamesOutput, layerNames[id-1])
	}

	return &VideoObjectDetection{&net, classes, layersNamesOutput}
}

// Perform object detection and return bounding boxes, confidences, and class IDs
func (vod *VideoObjectDetection) detectObjects(frame gocv.Mat) ([]image.Rectangle, []float32, []int) {
	width := frame.Cols()
	height := frame.Rows()

	// Get blob from the image
	blob := gocv.BlobFromImage(frame, 1.0/255, image.Point{416, 416},
		gocv.NewScalar(127.5, 0, 0, 0), true, false)

	// Set the input blob for YOLO
	vod.net.SetInput(blob, "")
	layerOutputs := vod.net.ForwardLayers(vod.layersNamesOutput)

	var boxes []image.Rectangle
	var confidences []float32
	var classIds []int

	// Iterate over the layer outputs
	for _, output := range layerOutputs {
		for i := 0; i < output.Rows(); i++ {
			data := output.Row(i)
			classId := -1
			confidence := float32(0)
			for j := 5; j < data.Cols(); j++ {
				if data.GetFloatAt(0, j) > confidence {
					classId = j - 5
					confidence = data.GetFloatAt(0, j)
				}
			}

			if confidence < 0.5 || classId == -1 {
				continue
			}

			centerX := int(data.GetFloatAt(0, 0) * float32(width))
			centerY := int(data.GetFloatAt(0, 1) * float32(height))
			w := int(data.GetFloatAt(0, 2) * float32(width))
			h := int(data.GetFloatAt(0, 3) * float32(height))

			xMin := centerX - w/2
			yMin := centerY - h/2
			box := image.Rectangle{Min: image.Point{xMin, yMin}, Max: image.Point{xMin + w, yMin + h}}
			boxes = append(boxes, box)
			confidences = append(confidences, confidence)
			classIds = append(classIds, classId)
		}
	}

	return boxes, confidences, classIds
}

func (vod *VideoObjectDetection) Execute(fileName, outfilePath string, showWindow bool) (*string, error) {

	// Open the video capture
	video, err := gocv.OpenVideoCapture(fileName)
	if err != nil {
		log.Fatalf("Error opening video file: %v", err)
		return nil, err
	}

	defer video.Close()

	// Create a window to display the results
	window := gocv.NewWindow("YOLO Object Detection")
	defer window.Close()

	// Get the frame width and height
	width := int(video.Get(gocv.VideoCaptureFrameWidth))
	height := int(video.Get(gocv.VideoCaptureFrameHeight))
	// Prepare to write the output video
	outputFilename := outfilePath + src.ChangeFileExtension(fileName, ".mp4")
	writer, err := gocv.VideoWriterFile(outputFilename, "mp4v", 25, width, height, true)
	if err != nil || !writer.IsOpened() {
		log.Fatalf("Error opening video writer")
		return nil, err
	}

	var colors []color.RGBA
	for i := 0; i < len(vod.classes); i++ {
		color := color.RGBA{uint8(rand.Uint32()) * 255, uint8(rand.Uint32()) * 255, uint8(rand.Uint32()) * 255, 0}
		colors = append(colors, color)
	}

	var frameCount int
	var totalTime float64

	// Process the video frame by frame
	for {
		frame := gocv.NewMat()
		if ok := video.Read(&frame); !ok {
			break
		}

		// Start measuring time for the frame
		start := time.Now()

		// Detect objects
		boxes, confidences, classIds := vod.detectObjects(frame)

		// Apply Non-Maximum Suppression
		indices := gocv.NMSBoxes(boxes, confidences, 0.5, 0.3)

		// Draw bounding boxes and labels on the frame
		for _, idx := range indices {
			classId := classIds[idx]
			class := vod.classes[classId]
			color := colors[classId]
			box := boxes[idx]
			label := fmt.Sprintf("%s: %.4f", class, confidences[idx])
			gocv.Rectangle(&frame, box, color, 2)
			gocv.PutText(&frame, label, image.Pt(box.Min.X, box.Min.Y-5), gocv.FontHersheySimplex, 0.5, color, 1)
		}

		// Write the processed frame to the output video
		writer.Write(frame)

		// Update frame count and total time
		frameCount++
		totalTime += time.Since(start).Seconds()
		fmt.Printf("Processed frame %d time: %.5f seconds\n", frameCount, totalTime)

		if showWindow {
			// Show the frame
			window.IMShow(frame)

			// Exit if 'Esc' key is pressed
			if window.WaitKey(1) == 27 {
				break
			}
		}
	}

	// Print final results
	fmt.Printf("Total frames processed: %d\n", frameCount)
	fmt.Printf("Total processing time: %.5f seconds\n", totalTime)
	fmt.Printf("FPS: %.1f\n", float64(frameCount)/totalTime)

	// Release resources
	video.Close()
	writer.Close()

	return &outputFilename, nil
}

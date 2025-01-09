package main

import (
	"fmt"
	"go-ocr/src"
	doc "go-ocr/src/documents"
	img "go-ocr/src/images"
	vid "go-ocr/src/videos"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Please provide the *Algorithm*, *Input file path* and *language* as an argument.")
		os.Exit(1)
	}

	algorithm := os.Args[1]

	// The first argument is the input file path
	inputFile := os.Args[2]

	language := os.Args[3]

	// Read the file content
	if ok := src.FileExists(inputFile); !ok {
		log.Fatalf("Error reading file: %v", inputFile)
	}

	switch algorithm {
	case "PLAIN_TEXT_EXTRACTION":
		{
			extractedText := doc.NewPlainTextExtractor().
				Execute(inputFile, language)
			if len(extractedText) == 0 {
				fmt.Printf("File: %s \nResult: No text extracted.\n", inputFile)
				break
			}
			fmt.Printf("File: %s \nResult: \n%s\n", inputFile, extractedText)
			break
		}

	case "HOCR_TEXT_EXTRACTION":
		{
			outfilePath, err := doc.NewHOCRTextExtractor("fonts/").
				Execute(inputFile, language, "output/generated-hocr/")
			if err != nil {
				fmt.Printf("File: %s \nResult: No text extracted.%s\n", inputFile, err)
				break
			}

			fmt.Printf("File: %s \nResult: \n%s\n", inputFile, *outfilePath)
			break
		}

	case "IMG_OBJECT_DETECTION":
		{
			detectedObjects := img.NewImageObjectDetector(
				"models/yolov3.weights",
				"models/yolov3.cfg",
				"models/coco.names").
				Execute(inputFile)

			if len(detectedObjects) == 0 {
				fmt.Printf("File: %s \nResult: No objects detected.\n", inputFile)
				break
			}

			fmt.Printf("File: %s\n", inputFile)
			for _, detectedObject := range detectedObjects {
				fmt.Printf("Result: %s with confidence: %2f.\n", detectedObject.ClassName, detectedObject.Confidence)
			}

			break
		}
	case "VIDEO_OBJECT_DETECTION":
		{
			outfilePath, err := vid.NewVideoObjectDetector(
				"models/yolov3.weights",
				"models/yolov3.cfg",
				"models/coco.names").
				Execute(inputFile, "output/generated-video/", true)

			if err != nil {
				fmt.Printf("File: %s \nResult: No objects detected.%s\n", inputFile, err)
				break
			}

			fmt.Printf("File: %s \nResult: \n%s\n", inputFile, *outfilePath)

			break
		}

	default:
		log.Fatal("Allowed algorithm are: 'PLAIN_TEXT_EXTRACTION', 'HOCR_TEXT_EXTRACTION', 'IMG_OBJECT_DETECTION','VIDEO_OBJECT_DETECTION'")
		os.Exit(1)
	}
}

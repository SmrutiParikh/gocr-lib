package main

import (
	"fmt"
	"go-ocr/src"
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
	case "TEXT_EXTRACTION":
		{
			extractedText := src.NewTextExtractionAlgorithm().
				Execute(inputFile, language)
			if len(extractedText) == 0 {
				fmt.Printf("File: %s \nResult: No text extracted.\n", inputFile)
				break
			}
			fmt.Printf("File: %s \nResult: \n%s\n", inputFile, extractedText)
			break
		}

	case "TEXT_EXTRACTION_HOCR":
		{
			outfilePath, err := src.NewTextExtractionHOCRAlgorithm("fonts/").
				Execute(inputFile, language, "hocr/")
			if err != nil {
				fmt.Printf("File: %s \nResult: No text extracted.\n", inputFile)
				break
			}

			fmt.Printf("File: %s \nResult: \n%s\n", inputFile, *outfilePath)
			break
		}

	case "OBJECT_DETECTION":
		{
			detectedObjects := src.NewObjectDetectionAlgorithm(
				"models/yolov3.weights",
				"models/yolov3.cfg",
				"models/coco.names").
				Execute(inputFile)
			if len(detectedObjects) == 0 {
				fmt.Printf("File: %s \nResult: No objects detected.\n", inputFile)
				break
			}

			for _, detectedObject := range detectedObjects {
				fmt.Printf("File: %s \nResult: %s with confidence: %2f.\n", inputFile, detectedObject.ClassName, detectedObject.Confidence)
			}

			break
		}

	default:
		log.Fatal("Allowed algorithm are: 'TEXT_EXTRACTION', 'TEXT_EXTRACTION_HOCR', 'OBJECT_DETECTION'")
		os.Exit(1)
	}
}

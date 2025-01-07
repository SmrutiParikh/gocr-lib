package src

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/otiai10/gosseract/v2"
	"github.com/signintech/gopdf"
)

type TextExtractionHOCRAlgorithm struct {
	tempFolder  string
	fontsFolder string
}

func NewTextExtractionHOCRAlgorithm(fontsFolder string) *TextExtractionHOCRAlgorithm {
	return &TextExtractionHOCRAlgorithm{"../temp/", fontsFolder}
}

func (h *TextExtractionHOCRAlgorithm) Execute(fileName, lang, outDir string) (*string, error) {
	if err := h.extractHocrText(fileName, lang); err != nil {
		return nil, err
	}

	text, boxes, pageWidth, pageHeight := h.extractTextAndBoundingBoxes("extracted-text.hocr")

	if len(text) == 0 || len(boxes) == 0 || pageWidth == 0 || pageHeight == 0 {
		return nil, fmt.Errorf("error extracting texts and boxes")
	}

	outFilePath := outDir + changeFileExtension(fileName, ".pdf")
	return &outFilePath, h.generatePDF(lang, outFilePath, pageWidth, pageHeight, text, boxes)
}

func (h *TextExtractionHOCRAlgorithm) extractHocrText(fileName, lang string) error {
	err := NewTextExtractionAlgorithm().preProcessImage(fileName)
	if err != nil {
		log.Fatal("Failed to preprocess image:", err)
		return err
	}

	// Use Tesseract to extract text in HOCR format
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage(lang)
	// Set the image from gocv into Tesseract
	err = client.SetImage(h.tempFolder + "processed-image.jpg")
	if err != nil {
		log.Fatal("Failed to set image to Tesseract:", err)
		return err
	}

	// Extract text in HOCR format
	hocrText, err := client.HOCRText()
	if err != nil {
		log.Fatalf("Error extracting HOCR: %v", err)
	}

	file, err := os.OpenFile(h.tempFolder+"extracted-text.hocr", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}

	defer file.Close()

	// Write the text to the file
	_, err = file.WriteString(hocrText)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

// Struct to hold the bounding box values for each word
type bbox struct {
	x1, y1, x2, y2 float64
}

func (h *TextExtractionHOCRAlgorithm) extractTextAndBoundingBoxes(hocrFilePath string) ([]string, []struct{ x1, y1, x2, y2 float64 }, float64, float64) {
	file, err := os.Open(h.tempFolder + hocrFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil, 0, 0
	}
	defer file.Close()

	// Parse the HOCR HTML file
	doc, err := html.Parse(file)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, nil, 0, 0
	}

	var text []string
	var boxes []struct{ x1, y1, x2, y2 float64 }
	var pageWidth, pageHeight float64

	var parseNode func(*html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			// Look for the title attribute containing page dimensions
			if isValidNode(*n, "ocr_page") {
				for _, attr := range n.Attr {
					if attr.Key == "title" {
						// Extract the page width and height from the title attribute
						// Format: "ocr_page 1 1 612 792"
						dims, err := getDimensions(attr)
						if err != nil {
							return
						}

						pageWidth = dims.x2
						pageHeight = dims.y2
					}
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "span" {
			// Get the "title" attribute that contains the bounding box information
			for _, attr := range n.Attr {
				if isValidNode(*n, "ocrx_word") {
					if attr.Key == "title" {
						// Parse the dims values
						dims, err := getDimensions(attr)
						if err != nil {
							return
						}

						// Standardize the bbox values
						standardizedBox := bbox{
							x1: dims.x1 / pageWidth,
							y1: dims.y1 / pageHeight,
							x2: dims.x2 / pageWidth,
							y2: dims.y2 / pageHeight,
						}

						boxes = append(boxes, standardizedBox)

					}
				}
			}
		} else if n.Type == html.TextNode {
			// Extract the text inside the span element
			data := strings.TrimSpace(n.Data)
			if len(data) > 0 && data != "\n" {
				text = append(text, n.Data)
			}
		}

		// Recurse through the nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}

	// Start parsing the root document node
	parseNode(doc)

	return text, boxes, pageWidth, pageHeight
}

func isValidNode(node html.Node, attrName string) bool {
	for _, attr := range node.Attr {
		if strings.HasPrefix(attr.Val, attrName) {
			return true
		}
	}

	return false
}

func getDimensions(attr html.Attribute) (bbox, error) {
	parts := strings.Split(attr.Val, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "bbox") {

			dimensions := strings.Split(part, " ")
			if len(dimensions) == 5 {
				x1, _ := strconv.Atoi(dimensions[1])
				y1, _ := strconv.Atoi(dimensions[2])
				x2, _ := strconv.Atoi(dimensions[3])
				y2, _ := strconv.Atoi(dimensions[4])

				return bbox{
					float64(x1),
					float64(y1),
					float64(x2),
					float64(y2),
				}, nil
			}
		}
	}

	return bbox{}, fmt.Errorf("bbox not found")
}

func (h *TextExtractionHOCRAlgorithm) generatePDF(lang, outputFilePath string, pageWidth, pageHeight float64,
	text []string, boxes []struct{ x1, y1, x2, y2 float64 }) error {
	// Initialize PDF
	const dpi = 72.0 // Assuming 72 DPI for simplicity (standard for many PDF libraries)

	// Convert from pixels to points (1 pixel = 1 point in this case as 1 inch = 72 points)
	pageWidthInPoints := pageWidth / dpi * 72
	pageHeightInPoints := pageHeight / dpi * 72

	pdf := gopdf.GoPdf{}
	config := gopdf.Config{PageSize: gopdf.Rect{W: pageWidthInPoints, H: pageHeightInPoints}}
	pdf.Start(config)
	pdf.AddPage()

	// Set font (make sure you have a font file or use a default font)
	if lang == "jpn" {
		err := pdf.AddTTFFont("Noto Sans", h.fontsFolder+"noto_sans.ttf")
		if err != nil {
			fmt.Println("Error loading font:", err)
			return err
		}
		err = pdf.SetFont("Noto Sans", "", 12)
		if err != nil {
			fmt.Println("Error setting font:", err)
			return err
		}
	} else {
		err := pdf.AddTTFFont("Arial", h.fontsFolder+"arial.ttf")
		if err != nil {
			fmt.Println("Error loading font:", err)
			return err
		}
		err = pdf.SetFont("Arial", "", 12)
		if err != nil {
			fmt.Println("Error setting font:", err)
			return err
		}
	}

	// Iterate through the text and bounding boxes and add text at the specified positions
	for i, t := range text {
		if i < len(boxes) {
			// Add the text at the bounding box position
			box := boxes[i]
			pdf.SetX(box.x1 * pageWidth)
			pdf.SetY(box.y1 * pageHeight)
			pdf.Cell(nil, t)
		}
	}

	// Write the output PDF
	err := pdf.WritePdf(outputFilePath)
	if err != nil {
		fmt.Println("Error writing PDF:", err)
		return err
	}

	return nil
}

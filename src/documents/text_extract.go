package doc

import (
	"image"
	"log"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type PlainTextExtractor struct {
	tempFolder string
}

func NewPlainTextExtractor() *PlainTextExtractor {
	return &PlainTextExtractor{"../../temp/"}
}

func (pte *PlainTextExtractor) Execute(fileName string, lang string) string {
	err := pte.preProcessImage(fileName)
	if err != nil {
		log.Fatal("Failed to preprocess image:", err)
		return err.Error()
	}

	// Now we will use Tesseract to extract text from the processed image
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage(lang)

	// Set the image to Tesseract
	err = client.SetImage(pte.tempFolder + "processed-image.jpg")
	if err != nil {
		log.Fatal("Failed to set image to Tesseract:", err)
		return err.Error()
	}

	// Extract text
	text, err := client.Text()
	if err != nil {
		log.Fatal("Failed to extract text:", err)
		return err.Error()
	}

	return text
}

func (pte *PlainTextExtractor) preProcessImage(fileName string) error {
	// Initialize ImageMagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Load the image with ImageMagick
	mw := imagick.NewMagickWand()

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(300, 300); err != nil {
		log.Fatal("Failed to set image resolution:", err)
		return err
	}

	err := mw.ReadImage(fileName)
	if err != nil {
		log.Fatal("Failed to read image:", err)
		return err
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
		log.Fatal("Failed to set alpha channel:", err)
		return err
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		log.Fatal("Unable to set compression quality:", err)
		return err
	}

	// Optionally, convert or process the image with ImageMagick if necessary.
	// Example: convert image to grayscale
	err = mw.SetImageColorspace(imagick.COLORSPACE_GRAY)
	if err != nil {
		log.Fatal("Failed to set colorspace:", err)
		return err
	}

	// Save the processed image
	err = mw.WriteImage(pte.tempFolder + "processed-image.jpg")
	if err != nil {
		log.Fatal("Failed to save processed image:", err)
		return err
	}

	// Load image using OpenCV for further processing
	img := gocv.IMRead(pte.tempFolder+"processed-image.jpg", gocv.IMReadColor)
	if img.Empty() {
		log.Fatal("Could not read the image")
		return nil
	}
	defer img.Close()

	// Optionally apply OpenCV transformations (e.g., resizing, thresholding)
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	gocv.GaussianBlur(img, &img, image.Pt(5, 5), 0, 0, gocv.BorderDefault)
	return nil
}

package images

// func main() {

// 	// Load the image
// 	img := gocv.IMRead("../samples/bill.jpg", gocv.IMReadColor)
// 	if img.Empty() {
// 		log.Fatal("Error loading image!")
// 	}
// 	defer img.Close()

// 	// Convert the image to grayscale
// 	gray := gocv.NewMat()
// 	defer gray.Close()
// 	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

// 	// Apply Gaussian blur to remove noise
// 	blurred := gocv.NewMat()
// 	defer blurred.Close()
// 	gocv.GaussianBlur(gray, &blurred, image.Pt(5, 5), 0)

// 	// Apply Canny edge detection to find edges
// 	edges := gocv.NewMat()
// 	defer edges.Close()
// 	gocv.Canny(blurred, &edges, 100, 200)

// 	// Find contours in the image
// 	contours := gocv.FindContours(edges, gocv.RetrievalExternal, gocv.ChainApproxSimple)

// 	// Draw contours and identify potential signatures
// 	for i, contour := range contours {
// 		// Only consider contours that are large enough (for example, signatures are typically larger)
// 		if len(contour) > 100 { // This threshold can be adjusted based on your image size
// 			// Draw contour on the original image
// 			gocv.DrawContours(&img, contours, i, color.RGBA{0, 255, 0, 0}, 2)

// 			// Optionally: analyze the bounding box or area to check if it matches signature characteristics
// 			boundingRect := gocv.BoundingRect(contour)
// 			fmt.Println("Bounding Rect:", boundingRect)
// 		}
// 	}

// 	// Show the result
// 	window := gocv.NewWindow("Signature Detection")
// 	defer window.Close()
// 	window.IMShow(img)
// 	window.WaitKey(0)
// }

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func main() {
	// 1. Open the Mac's built-in camera (Device 0)
	// On macOS, OpenCV automatically hooks into the AVFoundation backend
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalf("Error opening web cam: %v", err)
	}
	defer webcam.Close()

	// 2. Load the YOLO ONNX model file
	modelPath := "yolov26n.onnx" // Replace with your local YOLO26 model file path
	net := gocv.ReadNetFromONNX(modelPath)
	if net.Empty() {
		log.Fatalf("Error reading network from model file: %s", modelPath)
	}
	defer net.Close()

	// Set preferred backend and target (CPU is default and highly optimized on Apple Silicon)
	net.SetPreferableBackend(gocv.NetBackendDefault)
	net.SetPreferableTarget(gocv.NetTargetCPU)

	// Create a window to display the video feed
	window := gocv.NewWindow("GoCV YOLO26 Mac Camera")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	fmt.Println("Press 'q' in the camera window to exit.")

	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Device closed or unable to read from the webcam")
			break
		}

		// 3. Prepare the image frame as a Blob for the DNN network
		// YOLO models typically expect 640x640 input shapes
		blob := gocv.BlobFromImage(img, 1.0/255.0, image.Pt(640, 640), gocv.NewScalar(0, 0, 0, 0), true, false)
		net.SetInput(blob, "")

		// 4. Run forward pass to retrieve outputs
		outputs := net.Forward("")

		// 5. Parse bounding boxes and draw results
		// Note: YOLO26 is natively NMS-free, meaning outputs provide clean, raw final predictions
		processDetections(&img, outputs)

		// Close unused mats to prevent memory leaks
		blob.Close()
		outputs.Close()

		// Show the frame and check for exit keystroke
		window.IMShow(img)
		if window.WaitKey(1) == int('q') {
			break
		}
	}
}

// processDetections parses the network outputs and draws bounding boxes on the frame
func processDetections(frame *gocv.Mat, outputs gocv.Mat) {
	// YOLO output shape format typically varies by model variant.
	// For standard 2D detection, it contains coordinates [x_center, y_center, width, height, confidence...]
	// Loop over rows and isolate rows with a high confidence score:

	confidenceThreshold := float32(0.25)
	green := color.RGBA{0, 255, 0, 0}

	for i := 0; i < outputs.Rows(); i++ {
		confidence := outputs.GetFloatAt(i, 4)
		if confidence > confidenceThreshold {
			// Extract localized box geometry relative to 640x640 size
			centerX := outputs.GetFloatAt(i, 0)
			centerY := outputs.GetFloatAt(i, 1)
			width := outputs.GetFloatAt(i, 2)
			height := outputs.GetFloatAt(i, 3)

			// Scale the coordinates back up to the original frame canvas size
			scaleX := float32(frame.Cols()) / 640.0
			scaleY := float32(frame.Rows()) / 640.0

			left := int((centerX - width/2) * scaleX)
			top := int((centerY - height/2) * scaleY)
			right := int((centerX + width/2) * scaleX)
			bottom := int((centerY + height/2) * scaleY)

			// Render bounding box overlay
			rect := image.Rect(left, top, right, bottom)
			gocv.Rectangle(frame, rect, green, 2)

			// Optional label rendering setup
			gocv.PutText(frame, fmt.Sprintf("Object: %.2f", confidence), image.Pt(left, top-10),
				gocv.FontHersheySimplex, 0.5, green, 2)
		}
	}
}

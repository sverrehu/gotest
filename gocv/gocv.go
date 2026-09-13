package main

// On my Mac: PKG_CONFIG_PATH=/opt/local/lib/opencv4/pkgconfig go run gocv.go

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/born-ml/born/backend/cpu"
	"github.com/born-ml/born/onnx"
	"github.com/born-ml/born/tensor"
	"gocv.io/x/gocv"
)

func main() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalf("Error opening web cam: %v", err)
	}
	defer webcam.Close()

	be := cpu.New()
	modelPath := "yolo26_face_fp16.onnx"
	model, err := onnx.Load(modelPath, be)
	if err != nil {
		log.Fatalf("Error reading network from model file: %s: %v", modelPath, err)
	}

	// Create a window to display the video feed
	window := gocv.NewWindow("GoCV YOLO26 Mac Camera")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	resized := gocv.NewMat()
	defer resized.Close()

	rgbMat := gocv.NewMat()
	defer rgbMat.Close()

	inputBuffer := make([]float32, 3*640*640)

	fmt.Println("Press 'q' in the camera window to exit.")

	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Device closed or unable to read from the webcam")
			break
		}

		origH := img.Rows()
		origW := img.Cols()
		if origH == 0 || origW == 0 {
			continue
		}

		// Prepare 640x640 RGB normalized input for YOLO26
		gocv.Resize(img, &resized, image.Pt(640, 640), 0, 0, gocv.InterpolationLinear)
		gocv.CvtColor(resized, &rgbMat, gocv.ColorBGRToRGB)

		hwc, err := rgbMat.DataPtrUint8()
		if err != nil {
			log.Printf("Failed to get image bytes: %v", err)
			continue
		}

		for y := 0; y < 640; y++ {
			rowOffset := y * 640
			for x := 0; x < 640; x++ {
				srcIdx := (rowOffset + x) * 3
				inputBuffer[0*640*640+rowOffset+x] = float32(hwc[srcIdx+0]) / 255.0
				inputBuffer[1*640*640+rowOffset+x] = float32(hwc[srcIdx+1]) / 255.0
				inputBuffer[2*640*640+rowOffset+x] = float32(hwc[srcIdx+2]) / 255.0
			}
		}

		inputTensor, err := tensor.FromSlice(inputBuffer, tensor.Shape{1, 3, 640, 640}, be)
		if err != nil {
			log.Printf("Failed to create input tensor: %v", err)
			continue
		}

		// Run forward pass
		outputs, err := model.Forward(inputTensor.Raw())
		if err != nil {
			log.Printf("Forward pass error: %v", err)
			continue
		}

		// Parse bounding boxes and draw results
		processDetections(&img, outputs)

		// Show the frame and check for exit keystroke
		window.IMShow(img)
		if window.WaitKey(1) == int('q') {
			break
		}
	}
}

// processDetections parses the network outputs and draws bounding boxes on the frame
func processDetections(frame *gocv.Mat, outputs *tensor.RawTensor) {
	// YOLO26 output format: [batch=1, num_detections=300, 6]
	// Each detection row: [x1, y1, x2, y2, score, class]
	confidenceThreshold := float32(0.25)
	green := color.RGBA{0, 255, 0, 0}

	data := outputs.AsFloat32()
	numDetections := outputs.Shape()[1]
	rowLen := outputs.Shape()[2]

	scaleX := float32(frame.Cols()) / 640.0
	scaleY := float32(frame.Rows()) / 640.0

	for i := 0; i < numDetections; i++ {
		offset := i * rowLen
		x1 := data[offset+0]
		y1 := data[offset+1]
		x2 := data[offset+2]
		y2 := data[offset+3]
		score := data[offset+4]

		if score > confidenceThreshold {
			left := int(x1 * scaleX)
			top := int(y1 * scaleY)
			right := int(x2 * scaleX)
			bottom := int(y2 * scaleY)

			if left < 0 {
				left = 0
			}
			if top < 0 {
				top = 0
			}
			if right >= frame.Cols() {
				right = frame.Cols() - 1
			}
			if bottom >= frame.Rows() {
				bottom = frame.Rows() - 1
			}

			// Render bounding box overlay
			rect := image.Rect(left, top, right, bottom)
			gocv.Rectangle(frame, rect, green, 2)

			// Label rendering
			labelY := top - 10
			if labelY < 20 {
				labelY = top + 20
			}
			gocv.PutText(frame, fmt.Sprintf("Face: %.2f", score), image.Pt(left, labelY),
				gocv.FontHersheySimplex, 0.5, green, 2)
		}
	}
}

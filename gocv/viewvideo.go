package main

// On my Mac: PKG_CONFIG_PATH=/opt/local/lib/opencv4/pkgconfig go run viewvideo.go

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
		log.Fatalf("Error opening webcam: %v", err)
	}
	defer webcam.Close()

	be := cpu.New()
	modelPath := "yolo26_face_fp16.onnx"
	model, err := onnx.Load(modelPath, be)
	if err != nil {
		log.Fatalf("Error loading ONNX model %s: %v", modelPath, err)
	}

	window := gocv.NewWindow("YOLO26 Face Detection")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	resized := gocv.NewMat()
	defer resized.Close()

	rgbMat := gocv.NewMat()
	defer rgbMat.Close()

	inputBuffer := make([]float32, 3*640*640)
	boxColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	confidenceThreshold := float32(0.25)

	log.Println("Starting video stream with YOLO26 face detection. Press 'q' or ESC to exit.")

	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Device closed or unable to read from webcam")
			break
		}

		origH := img.Rows()
		origW := img.Cols()
		if origH == 0 || origW == 0 {
			continue
		}

		// Resize to 640x640 and convert BGR to RGB
		gocv.Resize(img, &resized, image.Pt(640, 640), 0, 0, gocv.InterpolationLinear)
		gocv.CvtColor(resized, &rgbMat, gocv.ColorBGRToRGB)

		// Convert HWC uint8 RGB to normalized CHW float32
		hwc, err := rgbMat.DataPtrUint8()
		if err != nil {
			log.Printf("Failed to get image bytes: %v", err)
			continue
		}
		for y := 0; y < 640; y++ {
			rowOffset := y * 640
			for x := 0; x < 640; x++ {
				srcIdx := (rowOffset + x) * 3
				r := float32(hwc[srcIdx+0]) / 255.0
				g := float32(hwc[srcIdx+1]) / 255.0
				b := float32(hwc[srcIdx+2]) / 255.0

				inputBuffer[0*640*640+rowOffset+x] = r
				inputBuffer[1*640*640+rowOffset+x] = g
				inputBuffer[2*640*640+rowOffset+x] = b
			}
		}

		inputTensor, err := tensor.FromSlice(inputBuffer, tensor.Shape{1, 3, 640, 640}, be)
		if err != nil {
			log.Printf("Failed to create input tensor: %v", err)
			continue
		}

		output, err := model.Forward(inputTensor.Raw())
		if err != nil {
			log.Printf("Forward pass error: %v", err)
			continue
		}

		// YOLO26 output format: [batch=1, num_detections=300, 6]
		// Each detection row: [x1, y1, x2, y2, score, class]
		data := output.AsFloat32()
		numDetections := output.Shape()[1]
		rowLen := output.Shape()[2]

		scaleX := float32(origW) / 640.0
		scaleY := float32(origH) / 640.0

		for i := 0; i < numDetections; i++ {
			offset := i * rowLen
			x1 := data[offset+0]
			y1 := data[offset+1]
			x2 := data[offset+2]
			y2 := data[offset+3]
			score := data[offset+4]

			if score >= confidenceThreshold {
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
				if right >= origW {
					right = origW - 1
				}
				if bottom >= origH {
					bottom = origH - 1
				}

				rect := image.Rect(left, top, right, bottom)
				gocv.Rectangle(&img, rect, boxColor, 2)

				label := fmt.Sprintf("Face: %.2f", score)
				labelY := top - 10
				if labelY < 20 {
					labelY = top + 20
				}
				gocv.PutText(&img, label, image.Pt(left, labelY),
					gocv.FontHersheySimplex, 0.5, boxColor, 2)
			}
		}

		window.IMShow(img)
		key := window.WaitKey(1)
		if key == int('q') || key == 27 {
			break
		}
	}
}

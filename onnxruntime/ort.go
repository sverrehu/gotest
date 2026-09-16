package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"

	"image/draw"
	_ "image/jpeg"

	ort "github.com/yalue/onnxruntime_go"
)

const onnxruntimeVersion = "1.29.0"
const onnxruntimeLibPath = "../../../lib/onnxruntime"

func findSharedLibrary() string {
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("%s/onnxruntime-osx-arm64-%s/lib/libonnxruntime.%s.dylib", onnxruntimeLibPath, onnxruntimeVersion, onnxruntimeVersion)
		}
	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("%s/onnxruntime-linux-aarch64-%s/lib/libonnxruntime.so.%s", onnxruntimeLibPath, onnxruntimeVersion, onnxruntimeVersion)
		}
	}
	log.Fatalf("Unable to determine a path to the onnxruntime shared library for OS \"%s\" and architecture \"%s\".\n",
		runtime.GOOS, runtime.GOARCH)
	return "//never gets here"
}

func main() {
	ort.SetSharedLibraryPath(findSharedLibrary())
	e := ort.InitializeEnvironment()
	if e != nil {
		panic(e)
	}
	defer ort.DestroyEnvironment()
	modelPath := "../gocv/yolo26_face_fp16.onnx"
	//modelPath := "../gocv/yolov8n-face.onnx"
	inputPath := "../born/faces.jpg"
	outputPath := "/tmp/faces_annotated.jpg"

	inputTensor, rgbaImg, err := loadAndProcessImage(inputPath)
	if err != nil {
		log.Panic(err)
	}
	defer inputTensor.Destroy()

	outputShape := ort.NewShape(1, 300, 6)
	outputData := make([]float32, 1*300*6)
	outputTensor, err := ort.NewTensor(outputShape, outputData)
	if err != nil {
		panic(err)
	}
	defer outputTensor.Destroy()

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"images"}, []string{"output0"}, // Tensor node names inside the ONNX model
		[]ort.ArbitraryTensor{inputTensor}, []ort.ArbitraryTensor{outputTensor}, nil)
	if err != nil {
		panic(err)
	}
	defer session.Destroy()

	err = session.Run()
	if err != nil {
		panic(err)
	}
	// YOLO26 output format: [batch=1, num_detections=300, 6]
	// Each detection row: [x1, y1, x2, y2, score, class]
	data := outputData
	numDetections := int(outputShape[1])
	rowLen := int(outputShape[2])

	confidenceThreshold := float32(0.25)
	boxColor := color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green bounding box
	detectedCount := 0

	for i := 0; i < numDetections; i++ {
		offset := i * rowLen
		x1 := data[offset+0]
		y1 := data[offset+1]
		x2 := data[offset+2]
		y2 := data[offset+3]
		score := data[offset+4]
		classID := int(data[offset+5])

		if score >= confidenceThreshold {
			detectedCount++
			log.Printf("Detection %d: class=%d score=%.4f bbox=[%.1f, %.1f, %.1f, %.1f]",
				detectedCount, classID, score, x1, y1, x2, y2)

			drawBoundingBox(rgbaImg, int(x1), int(y1), int(x2), int(y2), boxColor, 3)
		}
	}

	log.Printf("Augmented image with %d bounding boxes. Saving to %s...", detectedCount, outputPath)
	if err := saveJPEG(outputPath, rgbaImg); err != nil {
		log.Panic(fmt.Errorf("failed to save augmented image: %w", err))
	}
	log.Printf("Successfully saved augmented image to %s", outputPath)
}

func loadAndProcessImage(filePath string) (*ort.Tensor[float32], *image.RGBA, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, nil, err
	}

	bounds := src.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, src, bounds.Min, draw.Src)

	// Pre-allocate planar CHW float32 array (3 channels * 640 height * 640 width)
	data := make([]float32, 3*640*640)

	offsetX := 0
	offsetY := 0

	for y := 0; y < 640; y++ {
		for x := 0; x < 640; x++ {
			r, g, b, _ := rgbaImg.At(offsetX+x, offsetY+y).RGBA()

			idxR := 0*640*640 + y*640 + x
			idxG := 1*640*640 + y*640 + x
			idxB := 2*640*640 + y*640 + x

			data[idxR] = float32(r) / 65535.0
			data[idxG] = float32(g) / 65535.0
			data[idxB] = float32(b) / 65535.0
		}
	}

	inputShape := ort.NewShape(1, 3, 640, 640)
	t, err := ort.NewTensor(inputShape, data)
	if err != nil {
		return nil, nil, err
	}
	return t, rgbaImg, nil
}

func drawBoundingBox(img *image.RGBA, x1, y1, x2, y2 int, c color.Color, thickness int) {
	bounds := img.Bounds()

	// Clamp coordinates to image boundaries
	if x1 < bounds.Min.X {
		x1 = bounds.Min.X
	}
	if y1 < bounds.Min.Y {
		y1 = bounds.Min.Y
	}
	if x2 >= bounds.Max.X {
		x2 = bounds.Max.X - 1
	}
	if y2 >= bounds.Max.Y {
		y2 = bounds.Max.Y - 1
	}
	if x1 > x2 || y1 > y2 {
		return
	}

	// Horizontal lines
	for t := 0; t < thickness; t++ {
		for x := x1; x <= x2; x++ {
			if y1+t <= y2 {
				img.Set(x, y1+t, c)
			}
			if y2-t >= y1 {
				img.Set(x, y2-t, c)
			}
		}
	}

	// Vertical lines
	for t := 0; t < thickness; t++ {
		for y := y1; y <= y2; y++ {
			if x1+t <= x2 {
				img.Set(x1+t, y, c)
			}
			if x2-t >= x1 {
				img.Set(x2-t, y, c)
			}
		}
	}
}

// saveJPEG saves an image as JPEG with high quality.
func saveJPEG(filePath string, img image.Image) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

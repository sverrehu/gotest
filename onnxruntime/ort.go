package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"runtime"

	"image/draw"

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
	modelPath := "yolo26_face_fp16.onnx"
	//modelPath := "yolov8n-face.onnx"
	inputPath := "../born/faces.jpg"
	outputPath := "/tmp/faces_annotated.jpg"

	inputTensor, rgbaImg, err := loadAndProcessImage(inputPath)
	if err != nil {
		log.Panic(err)
	}
	defer inputTensor.Destroy()

	// Prepare empty output tensor shape allocation
	// Standard YOLO outputs are [1, 84, 8400] -> 84 channels (4 box coordinates + 80 classes), 8400 boxes
	outputShape := ort.NewShape(1, 84, 8400)
	outputData := make([]float32, 1*84*8400)
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

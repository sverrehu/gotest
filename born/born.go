package main

import (
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/born-ml/born/backend/cpu"
	"github.com/born-ml/born/onnx"
	"github.com/born-ml/born/tensor"
)

func main() {
	be := cpu.New()
	modelPath := "../gocv/yolo26_face_fp16.onnx"
	//modelPath := "../gocv/yolov8n-face.onnx"
	model, err := onnx.Load(modelPath, be)
	if err != nil {
		log.Panic(err)
	}
	inputTensor, err := processAndCropImage("faces.jpg", be)
	if err != nil {
		log.Panic(err)
	}
	output, err := model.Forward(inputTensor.Raw())
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Prediction execution completed successfully. Output tensor shape: %v", output.Shape())
}

// processAndCropImage opens an image, crops a 640x640 section from the top-left,
// and converts it into a CHW (Channels, Height, Width) tensor normalized to.
func processAndCropImage(filePath string, be tensor.Backend) (*tensor.Tensor[float32, tensor.Backend], error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Pre-allocate planar CHW float32 array (3 channels * 640 height * 640 width)
	data := make([]float32, 3*640*640)

	// Offsets for top-left cropping. Change these to adjust the crop origin.
	offsetX := 0
	offsetY := 0

	for y := 0; y < 640; y++ {
		for x := 0; x < 640; x++ {
			// Read pixel coordinate directly from the source image bounds
			r, g, b, _ := src.At(offsetX+x, offsetY+y).RGBA()

			// Map standard 0-65535 uint32 pixel values to normalized 0.0-1.0 float32
			idxR := 0*640*640 + y*640 + x
			idxG := 1*640*640 + y*640 + x
			idxB := 2*640*640 + y*640 + x

			data[idxR] = float32(r) / 65535.0
			data[idxG] = float32(g) / 65535.0
			data[idxB] = float32(b) / 65535.0
		}
	}

	// Wrap raw slice into a 4D tensor with shape [batch=1, channels=3, height=640, width=640]
	t, err := tensor.FromSlice(data, tensor.Shape{1, 3, 640, 640}, be)
	return t, nil
}

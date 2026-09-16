package main

import (
	"fmt"
	"log"
	"runtime"

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
}

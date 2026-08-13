package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	f, err := os.OpenFile("data.bin", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	numElements := 2
	fileSize := numElements * 4

	b, err := unix.Mmap(int(f.Fd()), 0, fileSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	defer unix.Munmap(b)

	int32s := unsafe.Slice((*int32)(unsafe.Pointer(&b[0])), numElements)

	fmt.Printf("int32s from file: %08x %08x\n", int32s[0], int32s[1])
}

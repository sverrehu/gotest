package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	fmt.Println("Endianness:", binary.NativeEndian.String())
}

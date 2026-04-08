package main

import (
	"github.com/sverrehu/gotest/versions/internal/webserver"
)

func main() {
	err := webserver.Run()
	if err != nil {
		panic(err)
	}
}

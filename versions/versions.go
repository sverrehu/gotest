package main

import (
	"fmt"

	"github.com/sverrehu/gotest/versions/internal/repos"
)

func main() {
	releases, err := maven.GetReleases("com.fasterxml.jackson.core", "jackson-core")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", releases)
}

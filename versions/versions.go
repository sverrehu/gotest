package main

import (
	"fmt"

	"github.com/sverrehu/gotest/versions/internal/repos"
)

func main() {
	query("com.fasterxml.jackson.core", "jackson-core")
	query("com.fasterxml.jackson.core", "QQQQ-core")
}

func query(groupId, artifactId string) {
	releases, err := maven.GetReleases(groupId, artifactId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", releases)
}

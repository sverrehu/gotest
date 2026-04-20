package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Data struct {
	Name    string
	AltName *string
	Map     map[string]string
}

type NewData struct {
	Name    string
	AltName *string
	Map     map[string]string
	Age     int
}

func main() {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(Data{Name: "John", Map: map[string]string{"key": "value"}})
	input := &NewData{}
	gob.NewDecoder(&buf).Decode(input)
	fmt.Println(*input)
}

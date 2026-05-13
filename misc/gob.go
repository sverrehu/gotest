package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
)

type Data struct {
	Name    string
	AltName *string
	Map     map[string]string
	mutex   sync.Mutex
}

type NewData struct {
	Name    string
	AltName *string
	Map     map[string]string
	Age     int
	mutex   sync.Mutex
}

func main() {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(Data{Name: "John", Map: map[string]string{"key": "value"}})
	input := &NewData{}
	gob.NewDecoder(&buf).Decode(input)
	fmt.Println(*input)
}

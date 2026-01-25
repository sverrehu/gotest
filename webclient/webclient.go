package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.vg.no", nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

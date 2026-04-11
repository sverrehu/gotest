package webclient

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/sverrehu/goutils/lrumap"
)

var cache *lrumap.LRUMap

func Get(url string) (string, error) {
	cached := cache.Get(url)
	if cached != nil {
		return cached.(string), nil
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	cache.Put(url, string(body))
	return string(body), nil
}

func init() {
	cache = lrumap.New(1000, time.Hour)
}

package maven

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sverrehu/gotest/versions/internal"
)

func GetReleases(groupId, artifactId string) ([]internal.Release, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://central.sonatype.com/solrsearch/select?wt=json&q=g:"+url.QueryEscape(groupId)+"+AND+a:"+url.QueryEscape(artifactId)+"&sort=v+desc", nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	return nil, nil
}

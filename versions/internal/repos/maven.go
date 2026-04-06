package maven

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sverrehu/gotest/versions/internal"
)

type fullResponse struct {
	Response struct {
		Docs []struct {
			V         string `json:"v"`
			Timestamp int64  `json:"timestamp"`
		} `json:"docs"`
	} `json:"response"`
}

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
	releases, err := translate(string(body))
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func translate(jsonResponse string) ([]internal.Release, error) {
	var resp fullResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, err
	}
	releases := make([]internal.Release, 0, len(resp.Response.Docs))
	for _, doc := range resp.Response.Docs {
		release := internal.Release{}
		release.Version = doc.V
		release.ReleasedAt = time.UnixMilli(doc.Timestamp)
		releases = append(releases, release)
	}
	return releases, nil
}

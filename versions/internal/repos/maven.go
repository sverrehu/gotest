package maven

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/sverrehu/gotest/versions/internal"
	"github.com/sverrehu/gotest/versions/internal/webclient"
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
	searchUrl := getSearchUrl(groupId, artifactId)
	body, err := webclient.Get(searchUrl)
	if err != nil {
		return nil, err
	}
	releases, err := translate(body)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func getSearchUrl(groupId, artifactId string) string {
	return "https://central.sonatype.com/solrsearch/select?wt=json&q=g:" + url.QueryEscape(groupId) + "+AND+a:" + url.QueryEscape(artifactId) + "&sort=v+desc"
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

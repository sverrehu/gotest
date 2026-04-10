package repos

// TODO: this will only fetch the 100 most recent releases.
// Sample: https://hub.docker.com/v2/repositories/library/ubuntu/tags?page=1&page_size=100

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/sverrehu/gotest/versions/internal"
	"github.com/sverrehu/gotest/versions/internal/webclient"
)

type OCIReleasesFetcher struct {
}

type fullOCIResponse struct {
	Count    int    `json:"count"`
	Next     any    `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Creator int `json:"creator"`
		ID      int `json:"id"`
		Images  []struct {
			Architecture string    `json:"architecture"`
			Features     string    `json:"features"`
			Variant      any       `json:"variant"`
			Digest       string    `json:"digest"`
			Os           string    `json:"os"`
			OsFeatures   string    `json:"os_features"`
			OsVersion    any       `json:"os_version"`
			Size         int       `json:"size"`
			Status       string    `json:"status"`
			LastPulled   time.Time `json:"last_pulled"`
			LastPushed   time.Time `json:"last_pushed"`
		} `json:"images"`
		LastUpdated         time.Time `json:"last_updated"`
		LastUpdater         int       `json:"last_updater"`
		LastUpdaterUsername string    `json:"last_updater_username"`
		Name                string    `json:"name"`
		Repository          int       `json:"repository"`
		FullSize            int       `json:"full_size"`
		V2                  bool      `json:"v2"`
		TagStatus           string    `json:"tag_status"`
		TagLastPulled       time.Time `json:"tag_last_pulled"`
		TagLastPushed       time.Time `json:"tag_last_pushed"`
		MediaType           string    `json:"media_type"`
		ContentType         string    `json:"content_type"`
		Digest              string    `json:"digest"`
	} `json:"results"`
}

func (rf OCIReleasesFetcher) GetReleases(pkg string) ([]internal.Release, error) {
	parts := regexp.MustCompile("[:/]").Split(pkg, -1)
	if len(parts) != 2 {
		return nil, &ReleasesFetcherError{Err: fmt.Errorf("expected two parts, separated by '/' in OCI package, got %s", pkg), IsParameterError: true}
	}
	return getOciReleases(parts[0], parts[1])
}

func getOciReleases(repo, image string) ([]internal.Release, error) {
	searchUrl := getOciSearchUrl(repo, image)
	body, err := webclient.Get(searchUrl)
	if err != nil {
		return nil, err
	}
	if body == "" {
		return make([]internal.Release, 0), nil
	}
	releases, err := translateOCIResponse(body)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func getOciSearchUrl(repo, image string) string {
	return "https://hub.docker.com/v2/repositories/" + url.PathEscape(repo) + "/" + url.PathEscape(image) + "/tags?page=1&page_size=100"
}

func translateOCIResponse(jsonResponse string) ([]internal.Release, error) {
	var resp fullOCIResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, err
	}
	releases := make([]internal.Release, 0, resp.Count)
	for _, result := range resp.Results {
		release := internal.Release{}
		release.Version = result.Name
		release.ReleasedAt = result.TagLastPushed
		releases = append(releases, release)
	}
	return releases, nil
}

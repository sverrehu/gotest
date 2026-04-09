package repos

import "github.com/sverrehu/gotest/versions/internal"

type ReleasesFetcher interface {
	GetReleases(pkg string) ([]internal.Release, error)
}

type ReleasesFetcherError struct {
	Err              error
	IsParameterError bool
}

func (e *ReleasesFetcherError) Error() string {
	return e.Err.Error()
}

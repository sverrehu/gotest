package main

import (
	"fmt"
)

type Fetcher interface {
	GetUrl() string
}

type FetcherBase struct {
	Fetcher
	firstPage int
}

type GitHubReleasesFetcher struct {
	FetcherBase
}

type GitHubTagsFetcher struct {
	FetcherBase
}

func NewFetcherBase(firstPage int) *FetcherBase {
	return &FetcherBase{
		firstPage: firstPage,
	}
}

func (f *FetcherBase) PrintUrl() {
	fmt.Println(f.GetUrl()) // this gives SIGSEGV
}

func NewGitHubReleasesFetcher() *GitHubReleasesFetcher {
	return &GitHubReleasesFetcher{
		FetcherBase: *NewFetcherBase(1),
	}
}

func NewGitHubTagsFetcher() *GitHubTagsFetcher {
	return &GitHubTagsFetcher{
		FetcherBase: *NewFetcherBase(2),
	}
}

func (rf *GitHubReleasesFetcher) GetUrl() string {
	return fmt.Sprintf("releases-%d", rf.firstPage)
}

func (rf *GitHubTagsFetcher) GetUrl() string {
	return fmt.Sprintf("tags-%d", rf.firstPage)
}

func PrintUrl(f Fetcher) {
	fmt.Printf("%s\n", f.GetUrl()) // this works
}

func main() {
	f1 := NewGitHubReleasesFetcher()
	//f1.PrintUrl()
	PrintUrl(f1)
	f2 := NewGitHubTagsFetcher()
	//f2.PrintUrl()
	PrintUrl(f2)
}

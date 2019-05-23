package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "41ef6dfed5be4bcb3e93f61313a72fd2c9546ad0"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user

	language := "rust"

	rs, err := getRepos(ctx, client, language)
	if err != nil {
		log.Fatalf("error when fetching repos, err: %q", err)
	}

	names := []string{}
	for _, r := range rs {
		names = append(names, *r.FullName)
	}

	file, err := json.MarshalIndent(names, "", " ")
	if err != nil {
		log.Fatalf("error when marshalling json, err: %q", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.json", language), file, 0644)
	if err != nil {
		log.Fatalf("error when writing json file, err: %q", err)
	}
}

func getRepos(ctx context.Context, c *github.Client, language string) ([]github.Repository, error) {
	var allRepos []github.Repository
	opt := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		if len(allRepos) >= 800 {
			return allRepos, nil
		}

		result, response, err := c.Search.Repositories(ctx, fmt.Sprintf("language:%s", language), opt)
		if err != nil {
			return nil, err
		}

		log.Printf("rate limit: %s", response.Rate)

		allRepos = append(allRepos, result.Repositories...)
		if response.NextPage == 0 {
			return allRepos, nil
		}

		opt.Page = response.NextPage
	}

}

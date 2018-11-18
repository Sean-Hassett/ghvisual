package main

import (
	ghv "./ghvisual"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

const configFile = "./config/config.json"

var config ghv.Configuration

type Commit struct {
	Author string
	Date   time.Time
	Size   int
}
type Repo struct {
	Name    string
	Forked  bool
	Commits []Commit
}

func main() {
	// open the configuration file
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	// create api client with oauth2 token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	commitsOpt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{Page: 1, PerPage: 100},
	}

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)

	var repoList []Repo

	i := 0
	for _, repo := range repos {
		var commits []Commit
		repoList = append(repoList, Repo{Name: *repo.Name, Forked: *repo.Fork, Commits: commits})
		for {
			commits, resp, err := client.Repositories.ListCommits(ctx, config.Username, *repo.Name, commitsOpt)
			if err != nil {
				fmt.Println(err)
			}
			for _, commit := range commits {
				if *commit.Commit.Author.Name == config.Username {
					repoList[i].Commits = append(repoList[i].Commits, Commit{Author: *commit.Commit.Author.Name, Date: *commit.Commit.Author.Date, Size: 0})
					if *repo.Name == "ghvisual" {
						fmt.Println(*commit.Commit.Message)
					}
				}
			}
			if resp.NextPage == 0 {
				break
			}
			commitsOpt.Page = resp.NextPage
		}
		i += 1
	}
	for _, repo := range repoList {
		fmt.Println()
		fmt.Println(repo.Name)
		for _, commit := range repo.Commits {
			fmt.Println(commit.Author)
			fmt.Println(commit.Date)
		}
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

const configFile = "ghvisual/config/config.json"

var config Configuration

type Configuration struct {
	Username string
	Token    string
}
type Commit struct {
	Author string
	Date   time.Time
	Size   int
}
type Repo struct {
	Name      string
	Owner     string
	Size      int
	Updated   github.Timestamp
	Language  string
	Commits   []Commit
}

func Retrieve() []Repo {
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
		if !*repo.Fork {
			if err != nil {
				fmt.Println(err)
			}
			l := "None"
			if repo.Language != nil {
				l = *repo.Language
			}

			repoList = append(repoList, Repo{Name: *repo.Name,
				Owner:     *repo.Owner.Login,
				Size:      *repo.Size,
				Updated:   *repo.UpdatedAt,
			    Language:  l})
			for {
				commits, resp, err := client.Repositories.ListCommits(ctx, config.Username, *repo.Name, commitsOpt)
				if err != nil {
					fmt.Println(err)
				}
				for _, commit := range commits {
					if commit.Committer != nil {
						repoList[i].Commits = append(repoList[i].Commits, Commit{Author: *commit.Commit.Author.Name, Date: *commit.Commit.Author.Date, Size: 0})
					}
				}
				if resp.NextPage == 0 {
					break
				}
				commitsOpt.Page = resp.NextPage
			}
			i += 1
		}
	}
	return repoList
}

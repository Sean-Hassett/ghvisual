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
)

const configFile = "./config/config.json"
const startAddr = "https://api.github.com/user"

type Links struct {
	Name string
	Url  string
}

var config ghv.Configuration

// Makes a HTTP GET request to the passed in URL. Parses JSON data from the body of the response
// and writes matching entries to a passed in struct.
func main() {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)

	for _, repo := range repos {
		if !*repo.Private {
			fmt.Println(*repo.Name)
		}
	}
}

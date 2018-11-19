package ghvisual

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const userAgent = "Sean-Hassett-ghvisual"

type Configuration struct {
	Username string
	Email    string
	Token    string
}
type Repos struct {
	Repos_url string
}
type Links struct {
	Name string
	Url  string
}

// Makes a HTTP GET request to the passed in URL. Parses JSON data from the body of the response
// and writes matching entries to a passed in struct.
func GetJson(config *Configuration, url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("User-Agent", userAgent)
	request.SetBasicAuth(config.Username, config.Token)

	response, err := myClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	ret, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return json.Unmarshal(ret, &target)
}

// Returns the API address for the repos of the user defined in config
func GetReposURL(config *Configuration, userAddr string) string {
	var data Repos
	GetJson(config, userAddr, &data)
	return data.Repos_url
}

// Returns a list of repos that the user defined in config is owner of, both the name of the repo
// and its API address
func GetReposList(config *Configuration, reposAddr string) []Links {
	var data []Links
	GetJson(config, reposAddr, &data)
	return data
}

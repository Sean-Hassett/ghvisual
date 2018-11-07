package tests

import (
	ghv "../ghvisual"
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestData struct {
	Current_user_url string
	Repos_url        string
}

const configFile = "../config/config.json"

// Test the response from an API call
func TestAPICall(t *testing.T) {
	rootAddr := "https://api.github.com"
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	config := ghv.Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	var data TestData
	ghv.GetJson(&config, rootAddr, &data)

	expected := "https://api.github.com/user"
	actual := data.Current_user_url
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got: '%s'", expected, actual)
	}
}

// Test getting a user's repo URL
func TestGetReposURL(t *testing.T) {
	userAddr := "https://api.github.com/user"
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	config := ghv.Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	expected := "https://api.github.com/users/" + config.Username + "/repos"
	actual := ghv.GetReposURL(&config, userAddr)
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got: '%s'", expected, actual)
	}
}

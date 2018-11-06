package tests

import (
	"../ghvisual"
	"encoding/json"
	"log"
	"os"
	"testing"
)

type Test struct {
	Current_user_url string
}

const configFile = "../ghvisual/config/config.json"

// Test the response from an API call
func TestAPICall(t *testing.T) {
	rootAddr := "https://api.github.com"
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	config := ghvisual.Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}
	data := Test{}
	ghvisual.GetJson(config, rootAddr, &data)

	expected := "https://api.github.com/user"
	actual := data.Current_user_url
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got: '%s'", expected, actual)
	}
}

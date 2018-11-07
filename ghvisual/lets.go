package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFile = "../ghvisual/config/config.json"
const userAgent = "Sean-Hassett_ghvisual"
const startAddr = "https://api.github.com/user"
const tempAddr = "https://api.github.com/users/Sean-Hassett/repos"

type Configuration struct {
	Username string
	Token    string
}

type Repos struct {
	Repos_url string
}
type Links struct {
	Name string
}

// Makes a HTTP GET request to the passed in URL. Parses JSON data from the body of the response
// and writes matching entries to a passed in struct.
func main() {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}
	result := GetJson(&config, tempAddr)
	var data []Links
	json.Unmarshal(result, &data)

	for _, name := range data {
		fmt.Println(name.Name)
	}
}

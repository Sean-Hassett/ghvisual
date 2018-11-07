package main

import (
	ghv "./ghvisual"
	"encoding/json"
	"fmt"
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
	linksAddr := ghv.GetReposURL(&config, startAddr)
	var data []Links
	ghv.GetJson(&config, linksAddr, &data)

	for _, name := range data {
		fmt.Println(name.Url)
	}
}

package ghvisual

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const configFile = "./config/config.json"
const userAgent = "Sean-Hassett_ghvisual"
const userAddr = "https://api.github.com/user"

type Configuration struct {
	Username string
	Token    string
}

type Json struct {
	Repos_url string
}

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

	data := Json{}
	GetJson(config, userAddr, &data)
	fmt.Println(data.Repos_url)
}

// Makes a HTTP GET request to the passed in URL. Parses JSON data from the body of the response
// and writes matching entries to a passed in struct.
func GetJson(config Configuration, url string, target interface{}) error {
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
	return json.NewDecoder(response.Body).Decode(target)
}

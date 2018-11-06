package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const userAgent = "Sean-Hassett_ghvisual"
const addr = "https://api.github.com"

type Json struct {
	Current_user_url string
	Starred_url      string
}

func main() {
	data := Json{}
	getJson(addr, &data)
	fmt.Println(data.Current_user_url)
	fmt.Println(data.Starred_url)
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("User-Agent", userAgent)

	response, err := myClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

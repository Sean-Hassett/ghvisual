package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"net/http"
)

const width = 1920
const height = 1080
const offset = 10
const bgShade = 220

func main() {
	http.Handle("/", http.HandlerFunc(draw))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func draw(w http.ResponseWriter, req *http.Request) {
	repoList := Retrieve()
	i := offset
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(bgShade, bgShade, bgShade))
	for _, repo := range repoList {
		s := "fill:green;stroke:none"
		canvas.Circle(i+int(math.Log(float64(repo.Size)))*10, height/2, int(math.Log(float64(repo.Size)))*10, s)
		i += ((int(math.Log(float64(repo.Size))) * 10) * 2) + offset
		fmt.Println(i)
	}
	canvas.End()
}

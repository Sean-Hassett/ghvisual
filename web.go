package main

import (
	"github.com/ajstarks/svgo"
	ghv "github.com/seanh/ghvisual/ghvisual"
	"log"
	"math"
	"net/http"
)

const width = 1920
const height = 1080
const offset = 25

func main() {
	http.Handle("/", http.HandlerFunc(draw))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func draw(w http.ResponseWriter, req *http.Request) {
	repoList := ghv.Retrieve()
	i := 0
	canvas := svg.New(w)
	canvas.Start(width, height)
	for _, repo := range repoList {
		canvas.Circle(i + int(math.Log(float64(repo.Size)))*10 + offset, height/2, int(math.Log(float64(repo.Size)))*10, "fill:none;stroke:black")
		i += int(math.Log(float64(repo.Size)))*10 + offset
	}
	canvas.End()
}

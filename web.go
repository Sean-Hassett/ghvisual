package main

import (
	"github.com/ajstarks/svgo"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(circ))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func circ(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(500, 500)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.Text(500/2, 500/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:black")
	s.End()
}

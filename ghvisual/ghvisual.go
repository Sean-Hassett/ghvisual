package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"net/http"
	"time"
)

const width = 1920
const height = 1080
const offset = 10
const bgShade = 220

func main() {
	http.Handle("/", http.HandlerFunc(draw))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// normalise the daysSinceUpdate value for each repo into the range specified by the upper and lower bounds
func normaliseDaysSinceUpdate(days []int, lower, upper int) []int {
	min, max := days[0], days[0]
	for _, day := range days {
		if day > max {
			max = day
		}
		if day < min {
			min = day
		}
	}

	var normalisedDays = []int{}
	for _, day := range days {
		normalisedDays = append(normalisedDays, lower + (day - min) * (upper-lower) / (max-min))
	}
	return normalisedDays
}

// thanks to https://stackoverflow.com/users/5987/mark-ransom for the python version of this function
func redistributeRGB(rgb []int) []int {
	threshold := 255
	var max int = rgb[0]
	if rgb[1] > max {
		max = rgb[1]
	}
	if rgb[2] > max {
		max = rgb[2]
	}

	if max <= threshold {
		return rgb
	}
	total := rgb[0] + rgb[1] + rgb[3]
	if total >= 3*threshold {
		return []int{threshold, threshold, threshold}
	}

	x := (3*threshold - total) / (3*max - total)
	gray := threshold - x*max
	return []int{gray + x*rgb[0], gray + x*rgb[0], gray + x*rgb[0]}
}

func draw(w http.ResponseWriter, req *http.Request) {
	repoList := Retrieve()
	i := offset
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(bgShade, bgShade, bgShade))
	deepSkyBlue := canvas.RGB(0,104,139)
	var daysSinceUpdate = []int{}
	for _, repo := range repoList {
		daysSinceUpdate = append(daysSinceUpdate, int(((time.Now().Sub(repo.Updated.Time)).Hours()) / 24))
	}
	normaliseDays := normaliseDaysSinceUpdate(daysSinceUpdate, 1, 100)
	fmt.Println(daysSinceUpdate)
	fmt.Println(normaliseDays)
	for _, repo := range repoList {
		s := deepSkyBlue
		canvas.Circle(i+int(math.Log(float64(repo.Size)))*10, height/2, int(math.Log(float64(repo.Size)))*10, s)
		i += ((int(math.Log(float64(repo.Size))) * 10) * 2) + offset
	}
	canvas.End()
}

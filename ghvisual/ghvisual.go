package main

import (
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const width = 1920
const height = 1080
const bgShade = 180
const buffer = 20
const degrees = 360
const degreeToRadian = 0.0174533

func main() {
	http.Handle("/", http.HandlerFunc(draw))
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// normalise a list of values into the range specified by the upper and lower bounds
func normaliseRange(list_in []int, lower, upper int) []int {
	min, max := list_in[0], list_in[0]
	for _, day := range list_in {
		if day > max {
			max = day
		}
		if day < min {
			min = day
		}
	}

	var normalisedDays = []int{}
	for _, day := range list_in {
		normalisedDays = append(normalisedDays, int(float64(lower+(day-min)*(upper-lower)/(max-min))))
	}
	return normalisedDays
}

// thanks to https://stackoverflow.com/users/5987/mark-ransom for the python version of this function
func redistributeRGB(rgb []int) []int {
	threshold := 255
	max := rgb[0]
	if rgb[1] > max {
		max = rgb[1]
	}
	if rgb[2] > max {
		max = rgb[2]
	}

	if max <= threshold {
		return rgb
	}
	total := rgb[0] + rgb[1] + rgb[2]
	if total >= 3*threshold {
		return []int{threshold, threshold, threshold}
	}

	x := (3*threshold - total) / (3*max - total)
	gray := threshold - x*max
	return []int{gray + x*rgb[0], gray + x*rgb[1], gray + x*rgb[2]}
}

func draw(w http.ResponseWriter, req *http.Request) {
	repoList := Retrieve()
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(bgShade, bgShade, bgShade))
	repoColor := []float64{127.0, 191.0, 191.0}
	userColor := []int{255, 192, 76}

	// map number of commits per day
	var activeDays = []int{0, 0, 0, 0, 0, 0, 0}
	Monday := 0
	Tuesday := 1
	Wednesday := 2
	Thursday := 3
	Friday := 4
	Saturday := 5
	Sunday := 6
	// map number of commits to time of day
	// Morning: 6:00am-11:59am
	// Afternoon: 12:00pm-5:59pm
	// Evening: 6:00pm-11:59pm
	// Night: 12:00am-5:59am
	var activeHours = []int{0, 0, 0, 0}
	Morning := 0
	Afternoon := 1
	Evening := 2
	Night := 3

	var daysSinceUpdate []int
	var diameters []float64
	sumOfDiameters := 0
	radius := 0.0
	for _, repo := range repoList {
		daysSinceUpdate = append(daysSinceUpdate, int(((time.Now().Sub(repo.Updated.Time)).Hours())/24))

		for _, commit := range repo.Commits {
			switch commit.Date.Weekday() {
			case time.Monday:
				activeDays[Monday] += 1
			case time.Tuesday:
				activeDays[Tuesday] += 1
			case time.Wednesday:
				activeDays[Wednesday] += 1
			case time.Thursday:
				activeDays[Thursday] += 1
			case time.Friday:
				activeDays[Friday] += 1
			case time.Saturday:
				activeDays[Saturday] += 1
			case time.Sunday:
				activeDays[Sunday] += 1
			}
			hour := commit.Date.Hour()
			switch {
			case hour < 6:
				activeHours[Night] += 1
			case hour < 12:
				activeHours[Morning] += 1
			case hour < 18:
				activeHours[Afternoon] += 1
			case hour < 24:
				activeHours[Evening] += 1
			}
		}
		diameter := math.Log(float64(repo.Size))*10*2
		diameters = append(diameters, diameter)
		sumOfDiameters += int(diameter + buffer)
	}

	normalisedDays := normaliseRange(daysSinceUpdate, 5, 60)
	radius = float64(sumOfDiameters) / (2.0 * math.Pi)
	theta := -(buffer / float64(sumOfDiameters) * degrees)
	minX := width

	for i, repo := range repoList {
		brightnessMult := float64(normalisedDays[i]) / 60.0
		circleColor := redistributeRGB([]int{int(repoColor[0] * brightnessMult), int(repoColor[1] * brightnessMult), int(repoColor[2] * brightnessMult)})
		s := canvas.RGB(int(circleColor[0]), int(circleColor[1]), int(circleColor[2]))
		//angle in radians
		offset := buffer / float64(sumOfDiameters) * degrees
		theta += offset/2 + (diameters[i] / 2 / float64(sumOfDiameters) * degrees)
		angle := theta * degreeToRadian

		xVal := int(radius*math.Cos(angle)+width/2)
		if xVal - int(diameters[i] / 2) < minX {
			minX = xVal- int(diameters[i] / 2)
		}
		yVal := int(radius*math.Sin(angle)+height/2)
		canvas.Line(xVal, yVal, width/2, height/2, "stroke:black")
		canvas.Circle(xVal, yVal, int((math.Log(float64(repo.Size)))*10), s)
		canvas.Text(xVal, yVal+4, strconv.Itoa(i+1), "fill:white;text-anchor:middle")

		theta += offset/2 + (diameters[i] / 2 / float64(sumOfDiameters) * degrees)

	}
	canvas.Circle(width/2, height/2, int(radius/2), canvas.RGB(userColor[0], userColor[1], userColor[2]))
	canvas.Text(width/2, height/2+20, config.Username, "fill:black;text-anchor:middle;font-size:40px")

	chartBuffer := 50
	chartXStart := chartBuffer
	chartYStart := chartBuffer
	chartWidth := minX - chartBuffer*2
	chartHeight := height/2 - (chartBuffer)
	space := 10

	dayBarWidth := (chartWidth-space) / 7 - space
	canvas.Rect(chartXStart, chartYStart, chartWidth, chartHeight, "fill:rgb(140,140,140)")
	activeDaysValues := []int{}
	for _, day := range activeDays{
		activeDaysValues = append(activeDaysValues, day)
	}
	normalActiveDaysValues := normaliseRange(activeDaysValues, 100, chartHeight-chartBuffer)
	for i := range normalActiveDaysValues {
		p := ""
		canvas.Rect(chartXStart+space+(dayBarWidth + space)*i, chartYStart+chartHeight-normalActiveDaysValues[i], dayBarWidth, normalActiveDaysValues[i], "fill:rgb(51,78,78")
	}
	canvas.End()
}

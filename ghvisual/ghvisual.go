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
const degreeToRadian = 0.0174533

func main() {
	http.Handle("/", http.HandlerFunc(draw))
	err := http.ListenAndServe(":8082", nil)
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
		normalisedDays = append(normalisedDays, int(lower+(day-min)*(upper-lower)/(max-min)))
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
	deepSkyBlue := []float64{126.0, 192.0, 238.0}
	orange := []int{255, 199, 92}

	// map number of commits per day
	var activeDays = map[string]int{
		"Monday":    0,
		"Tuesday":   0,
		"Wednesday": 0,
		"Thursday":  0,
		"Friday":    0,
		"Saturday":  0,
		"Sunday":    0,
	}
	// map number of commits to time of day
	// Morning: 6:00am-11:59am
	// Afternoon: 12:00pm-5:59pm
	// Evening: 6:00pm-11:59pm
	// Night: 12:00am-5:59am
	var activeHours = map[string]int{
		"Morning":   0,
		"Afternoon": 0,
		"Evening":   0,
		"Night":     0,
	}

	var daysSinceUpdate []int
	var diameters []float64
	sumOfDiameters := 0
	radius := 0.0
	for _, repo := range repoList {
		daysSinceUpdate = append(daysSinceUpdate, int(((time.Now().Sub(repo.Updated.Time)).Hours())/24))

		for _, commit := range repo.Commits {
			switch commit.Date.Weekday() {
			case time.Monday:
				activeDays["Monday"] += 1
			case time.Tuesday:
				activeDays["Tuesday"] += 1
			case time.Wednesday:
				activeDays["Wednesday"] += 1
			case time.Thursday:
				activeDays["Thursday"] += 1
			case time.Friday:
				activeDays["Friday"] += 1
			case time.Saturday:
				activeDays["Saturday"] += 1
			case time.Sunday:
				activeDays["Sunday"] += 1
			}
			hour := commit.Date.Hour()
			switch {
			case hour < 6:
				activeHours["Night"] += 1
			case hour < 12:
				activeHours["Morning"] += 1
			case hour < 18:
				activeHours["Afternoon"] += 1
			case hour < 24:
				activeHours["Evening"] += 1
			}
		}
		diameter := math.Log(float64(repo.Size))*10*2
		diameters = append(diameters, diameter)
		sumOfDiameters += int(diameter + buffer)
	}

	normalisedDays := normaliseDaysSinceUpdate(daysSinceUpdate, 1, 100)
	radius = float64(sumOfDiameters) / (2.0 * math.Pi)
	theta := -(buffer / float64(sumOfDiameters) * 360)

	for i, repo := range repoList {
		brightnessMult := float64(normalisedDays[i]) / 100.0
		circleColor := redistributeRGB([]int{int(deepSkyBlue[0] * brightnessMult), int(deepSkyBlue[1] * brightnessMult), int(deepSkyBlue[2] * brightnessMult)})
		s := canvas.RGB(int(circleColor[0]), int(circleColor[1]), int(circleColor[2]))

		//angle in radians
		offset := buffer / float64(sumOfDiameters) * 360
		theta += offset/2 + (diameters[i] / 2 / float64(sumOfDiameters) * 360)
		angle := theta * degreeToRadian

		xVal := int(radius*math.Cos(angle)+width/2)
		yVal := int(radius*math.Sin(angle)+height/2)
		canvas.Line(xVal, yVal, width/2, height/2, "fill:black:weight:2")
		canvas.Circle(xVal, yVal, int((math.Log(float64(repo.Size)))*10), s)
		canvas.Text(xVal, yVal+4, strconv.Itoa(i+1), "fill:white;text-anchor:middle")
		canvas.Line(123, 123, width/2, height/2, "fill:black:weight:2")

		theta += offset/2 + (diameters[i] / 2 / float64(sumOfDiameters) * 360)

	}
	canvas.Circle(width/2, height/2, int(radius/2), canvas.RGB(orange[0], orange[1], orange[2]))
	canvas.Text(width/2, height/2+20, config.Username, "fill:black;text-anchor:middle;font-size:40px")
	canvas.End()
}

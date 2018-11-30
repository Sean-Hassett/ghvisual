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
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// normalise a list of values into the range specified by the upper and lower bounds
func normaliseRange(listIn []int, lower, upper int) []int {
	min, max := listIn[0], listIn[0]
	for _, day := range listIn {
		if day > max {
			max = day
		}
		if day < min {
			min = day
		}
	}

	var normalisedDays []int
	for _, day := range listIn {
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
	repoColor := []float64{79.0, 125.0, 125.0}
	userColor := []int{253, 183, 125}

	// map number of commits per day
	var days = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
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
	var hours = []string{"Morning", "Afternoon", "Evening", "Night"}
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
		diameter := math.Log(float64(repo.Size)) * 10 * 2
		diameters = append(diameters, diameter)
		sumOfDiameters += int(diameter + buffer)
	}

	normalisedDays := normaliseRange(daysSinceUpdate, 30, 90)
	radius = float64(sumOfDiameters) / (2.0 * math.Pi)
	if radius < 300 {
		radius = 300.0
	}
	theta := -(buffer / float64(sumOfDiameters) * degrees)
	minX := width
	maxX := 0

	for i, repo := range repoList {
		brightnessMult := float64(normalisedDays[i]) / 60.0
		circleColor := redistributeRGB([]int{int(repoColor[0] * brightnessMult), int(repoColor[1] * brightnessMult), int(repoColor[2] * brightnessMult)})
		s := canvas.RGB(int(circleColor[0]), int(circleColor[1]), int(circleColor[2]))
		//angle in radians
		offset := buffer / float64(sumOfDiameters) * degrees
		theta += offset/2 + (diameters[i] / 2 / float64(sumOfDiameters) * degrees)
		angle := theta * degreeToRadian

		xVal := int(radius*math.Cos(angle) + width/2)
		if xVal-int(diameters[i]/2) < minX {
			minX = xVal - int(diameters[i]/2)
		}
		if xVal+int(diameters[i]/2) > maxX {
			maxX = xVal + int(diameters[i]/2)
		}
		yVal := int(radius*math.Sin(angle) + height/2)
		canvas.Line(xVal, yVal, width/2, height/2, "stroke:black")
		canvas.Circle(xVal, yVal, int((math.Log(float64(repo.Size)))*10), s)
		canvas.Text(xVal, yVal+4, strconv.Itoa(i+1), "fill:rgb(220,220,220);text-anchor:middle")

		theta += offset/2 + (diameters[i] / 2 / float64(sumOfDiameters) * degrees)

	}
	canvas.Circle(width/2, height/2, int(radius/2), canvas.RGB(userColor[0], userColor[1], userColor[2]))
	canvas.Text(width/2, height/2+20, config.Username, "fill:rgb(23,31,31);text-anchor:middle;font-size:40px")

	// daily activity chart
	chartBuffer := 50
	chartXStart := chartBuffer
	chartYStart := chartBuffer
	chartWidth := minX - chartBuffer*2
	chartHeight := height/2 - (chartBuffer * 3 / 2)
	space := 10

	chartBackgroundStyle := "fill:rgb(140,140,140)"
	chartTextStyle := "fill:rgb(220,220,220);text-anchor:middle;font-size:18"
	barStyle := "fill:rgb(51,78,78)"
	barTextStyle := "text-anchor:middle;fill:rgb(220,220,220)"

	dayBarWidth := (chartWidth-space)/7 - space
	canvas.Rect(chartXStart, chartYStart, chartWidth, chartHeight, chartBackgroundStyle)
	canvas.Text(chartXStart+chartWidth/2, chartYStart+20, "Commits per Day of Week", chartTextStyle)
	var activeDaysValues []int
	for _, day := range activeDays {
		activeDaysValues = append(activeDaysValues, day)
	}
	normalActiveDaysValues := normaliseRange(activeDaysValues, 100, chartHeight-chartBuffer)
	for i := range normalActiveDaysValues {
		numCommits := strconv.Itoa(activeDaysValues[i])
		day := days[i]
		canvas.Rect(chartXStart+space+(dayBarWidth+space)*i, chartYStart+chartHeight-normalActiveDaysValues[i], dayBarWidth, normalActiveDaysValues[i], barStyle)
		canvas.Text(chartXStart+space+(dayBarWidth+space)*i+dayBarWidth/2, chartYStart+chartHeight-normalActiveDaysValues[i]+normalActiveDaysValues[i]-10, day, barTextStyle)
		canvas.Text(chartXStart+space+(dayBarWidth+space)*i+dayBarWidth/2, chartYStart+chartHeight-normalActiveDaysValues[i]+20, numCommits, barTextStyle)
	}

	//hourly activity chart
	hourBarWidth := (chartWidth-space)/4 - space
	chartYStart = chartYStart + chartHeight + chartBuffer
	canvas.Rect(chartXStart, chartYStart, chartWidth, chartHeight, chartBackgroundStyle)
	canvas.Text(chartXStart+chartWidth/2, chartYStart+20, "Commits per Time of Day", chartTextStyle)
	var activeHoursValues []int
	for _, day := range activeHours {
		activeHoursValues = append(activeHoursValues, day)
	}
	normalActiveHoursValues := normaliseRange(activeHoursValues, 100, chartHeight-chartBuffer)
	for i := range normalActiveHoursValues {
		numCommits := strconv.Itoa(activeHoursValues[i])
		hour := hours[i]
		canvas.Rect(chartXStart+space+(hourBarWidth+space)*i, chartYStart+chartHeight-normalActiveHoursValues[i], hourBarWidth, normalActiveHoursValues[i], barStyle)
		canvas.Text(chartXStart+space+(hourBarWidth+space)*i+hourBarWidth/2, chartYStart+chartHeight-normalActiveHoursValues[i]+normalActiveHoursValues[i]-10, hour, barTextStyle)
		canvas.Text(chartXStart+space+(hourBarWidth+space)*i+hourBarWidth/2, chartYStart+chartHeight-normalActiveHoursValues[i]+20, numCommits, barTextStyle)
	}

	//legend
	legendBuffer := 50
	legendXStart := maxX + legendBuffer
	legendYStart := chartBuffer
	legendWidth := width - legendXStart - legendBuffer
	legendHeight := height - legendBuffer*2

	legendTextYStart := legendBuffer * 3
	legendTextXStart := legendXStart + 10
	indexWidth := 30
	languageWidth := 90
	nameWidth := legendWidth - indexWidth - languageWidth - 20
	lineHeight := 40
	canvas.Rect(legendXStart, legendYStart, legendWidth, legendHeight, chartBackgroundStyle)

	legendTextStyle := "fill:rgb(220,220,220);text-anchor:middle;font-size:30"
	legendLineTextStyle := "fill:rgb(51,78,78);font-size:20"

	canvas.Text(legendXStart+legendWidth/2, legendYStart+36, "Legend", legendTextStyle)
	canvas.Text(legendTextXStart+indexWidth, legendTextYStart, "Name", legendLineTextStyle)
	canvas.Text(legendTextXStart+indexWidth+nameWidth, legendTextYStart, "Language", legendLineTextStyle)

	for i, repo := range repoList {
		strI := strconv.Itoa(i + 1)
		canvas.Text(legendTextXStart, legendTextYStart+(lineHeight*(i+1)), strI, legendLineTextStyle)
		canvas.Text(legendTextXStart+indexWidth, legendTextYStart+(lineHeight*(i+1)), repo.Name, legendLineTextStyle)
		canvas.Text(legendTextXStart+indexWidth+nameWidth, legendTextYStart+(lineHeight*(i+1)), repo.Language, legendLineTextStyle)
	}
	canvas.End()
}

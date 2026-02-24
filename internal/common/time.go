package common

import (
	"math"
	"os"
	"time"
)

func GetClosestQuarterHour() time.Time {
	now := time.Now()
	currentHour := now.Hour()
	currentMinutes := now.Minute()

	if currentMinutes > 52 {
		currentHour++
	}

	location, _ := LoadLocalTime()

	m1 := math.Round(float64(currentMinutes)/float64(15)) * 15
	m2 := int(m1) % 60

	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		currentHour,
		m2,
		0,
		0,
		location,
	)
}

func LoadLocalTime() (*time.Location, error) {
	l := os.Getenv("TZ")
	location, err := time.LoadLocation(l)

	if err != nil {
		return nil, err
	}

	return location, nil
}

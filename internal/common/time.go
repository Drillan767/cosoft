package common

import (
	"math"
	"time"
)

func GetClosestQuarterHour() time.Time {
	now := time.Now()
	currentHour := now.Hour()
	currentMinutes := now.Minute()

	if currentMinutes > 52 {
		currentHour++
	}

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
		time.UTC,
	)
}

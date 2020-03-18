package apod

import (
	"math/rand"
	"time"
)

// RandDate returns a random date between 1995-06-16 & today
func RandDate() string {
	var (
		min   int64
		max   int64
		sec   int64
		delta int64
		date  string
	)
	min = time.Date(1995, 6, 16, 0, 0, 0, 0, time.UTC).Unix()
	max = time.Now().UTC().Unix()
	delta = max - min

	sec = rand.Int63n(delta) + min
	date = time.Unix(sec, 0).Format("2006-01-02")

	return date
}

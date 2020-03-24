package apod

import (
	"math/rand"
	"time"
)

// RandDate returns a random date between 1995-06-16 & today.
func RandDate() string {
	var (
		min   int64
		max   int64
		sec   int64
		delta int64
		date  string
	)
	min = time.Date(1995, 6, 16, 0, 0, 0, 0, time.UTC).Unix()

	// We are taking max from UTC but it could fail if the
	// timezone on NASA APOD server is any different. They don't
	// mention the timezone api server uses so we use UTC & the
	// probability of this function failing because of this issue
	// is very low.
	max = time.Now().UTC().Unix()
	delta = max - min

	sec = rand.Int63n(delta) + min
	date = time.Unix(sec, 0).Format("2006-01-02")

	return date
}

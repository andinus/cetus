package nasa

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"framagit.org/andinus/cetus/pkg/cetus"
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

// GetApodJson returns json response received from the api
func GetApodJson(reqInfo map[string]string, t time.Duration) (string, error) {
	re := regexp.MustCompile("((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])")
	if !re.MatchString(reqInfo["date"]) {
		return "", fmt.Errorf("%s does not match format 'YYYY-MM-DD'", reqInfo["date"])
	}

	params := make(map[string]string)
	params["api_key"] = reqInfo["apiKey"]
	params["date"] = reqInfo["date"]

	body, err := cetus.GetRes(reqInfo["api"], params, t)
	return string(body), err
}

package nasa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

// APOD holds responses
type APOD struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HDURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`

	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

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

// APODPath returns Astronomy Picture of the Day path
func APODPath(apodInfo map[string]string, timeout time.Duration) (APOD, error) {
	var err error
	apodRes := APOD{}

	// validate date
	re := regexp.MustCompile("((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])")
	if !re.MatchString(apodInfo["date"]) {
		return apodRes, fmt.Errorf("%s does not match format 'YYYY-MM-DD'", apodInfo["date"])
	}

	client := http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest(http.MethodGet, apodInfo["api"], nil)
	if err != nil {
		return apodRes, err
	}
	q := req.URL.Query()
	q.Add("api_key", apodInfo["apiKey"])
	q.Add("date", apodInfo["date"])
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: GET %s\n", apodInfo["api"])
		return apodRes, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return apodRes, err
	}

	err = json.Unmarshal([]byte(resBody), &apodRes)
	if err != nil {
		return apodRes, err
	}

	if res.StatusCode != 200 {
		return apodRes, fmt.Errorf("Unexpected response status code received: %d %s",
			res.StatusCode, http.StatusText(res.StatusCode))
	}

	return apodRes, err
}

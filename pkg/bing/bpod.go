package bing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Photo holds responses
type Photo struct {
	StartDate     string `json:"startdate"`
	FullStartDate string `json:"fullstartdate"`
	EndDate       string `json:"enddate"`
	URL           string `json:"url"`
	URLBase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Hsh           string `json:"hsh"`
}

// BPOD  holds list of response
type BPOD struct {
	Photos []Photo `json:"images"`
}

// BPODPath returns Bing Photo of the Day responses
func BPODPath(bpodInfo map[string]string, timeout time.Duration) (BPOD, error) {
	var err error
	bpodRes := BPOD{}

	client := http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest(http.MethodGet, bpodInfo["api"], nil)
	if err != nil {
		return bpodRes, err
	}
	q := req.URL.Query()
	q.Add("format", "js")

	// if random flag is passed then fetch 7 photos
	if bpodInfo["random"] == "true" {
		q.Add("n", "7")
	} else {
		q.Add("n", "1")
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		return bpodRes, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return bpodRes, err
	}

	err = json.Unmarshal([]byte(resBody), &bpodRes)
	if err != nil {
		return bpodRes, err
	}

	if res.StatusCode != 200 {
		return bpodRes, fmt.Errorf("Unexpected response status code received: %d %s",
			res.StatusCode, http.StatusText(res.StatusCode))
	}

	return bpodRes, err
}

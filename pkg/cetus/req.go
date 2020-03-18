package cetus

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetRes returns api response
func GetRes(api string, params map[string]string) (string, error) {
	c := http.Client{
		// TODO: timeout should be configurable by the user
		Timeout: time.Second * 64,
	}

	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Unexpected response status code received: %d %s",
			res.StatusCode, http.StatusText(res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

package apod

import (
	"encoding/json"
	"fmt"
	"regexp"

	"framagit.org/andinus/cetus/pkg/request"
)

// Res holds the response from the api.
type Res struct {
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

// UnmarshalJson will take body as input & unmarshal it to res
func UnmarshalJson(res *Res, body string) error {
	err := json.Unmarshal([]byte(body), res)
	if err != nil {
		return fmt.Errorf("UnmarshalJson failed\n%s", err.Error())
	}
	return nil
}

// GetJson returns json response received from the api
func GetJson(reqInfo map[string]string) (string, error) {
	re := regexp.MustCompile("((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])")
	if !re.MatchString(reqInfo["date"]) {
		return "", fmt.Errorf("%s does not match format 'YYYY-MM-DD'", reqInfo["date"])
	}

	params := make(map[string]string)
	params["api_key"] = reqInfo["apiKey"]
	params["date"] = reqInfo["date"]

	body, err := request.GetRes(reqInfo["api"], params)
	return string(body), err
}

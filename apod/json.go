package apod

import (
	"encoding/json"
	"fmt"
	"regexp"

	"framagit.org/andinus/cetus/pkg/request"
)

// APOD holds the response from the api. Not every field is returned
// in every request. Code & Msg should be filled only if the api
// returns an error, this behaviour was observed and shouldn't be
// trusted.
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

// UnmarshalJson will take body as input & unmarshal it to res.
func UnmarshalJson(res *Res, body string) error {
	err := json.Unmarshal([]byte(body), res)
	if err != nil {
		err = fmt.Errorf("json.go: unmarshalling json failed\n%s",
			err.Error())
	}
	return err
}

// GetJson takes reqInfo as input and returns the body and an error.
func GetJson(reqInfo map[string]string) (string, error) {
	var body string
	var err error

	// This regexp is not perfect and does not guarantee that the
	// request will not fail because of wrong date, this will
	// eliminate many wrong dates though.
	re := regexp.MustCompile("((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])")
	if !re.MatchString(reqInfo["date"]) {
		err = fmt.Errorf("json.go: %s does not match format 'YYYY-MM-DD'",
			reqInfo["date"])
		return body, err
	}

	// reqInfo is map[string]string and params is built from it, currently
	// it takes apiKey and the date from reqInfo to build param. If any
	// new key/value is added to reqInfo then it must be addded here too,
	// it won't be sent as param directly.
	params := make(map[string]string)
	params["api_key"] = reqInfo["apiKey"]
	params["date"] = reqInfo["date"]

	body, err = request.GetRes(reqInfo["api"], params)
	return body, err
}

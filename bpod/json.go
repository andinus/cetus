package bpod

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"tildegit.org/andinus/cetus/request"
)

// BPOD holds the response from the api.
type BPOD struct {
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

// List holds the list of BPOD.
type List struct {
	Photos []BPOD `json:"images"`
}

// MarshalJson takes res as input and returns body. This remarshaling
// is required because of a bug. To learn about why this is required,
// remove this function & then run `cetus set bpod -random`. Put a
// `fmt.Println(res, body)` somewhere and look at how they differ. res
// will contain a single entry but body will have all 7 entries which
// is bad because body is cached to disk to view later.
func MarshalJson(res BPOD) (string, error) {
	out, err := json.Marshal(res)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"MarshalJson failed",
			err.Error())
	}
	body := string(out)
	return body, err
}

// UnmarshalJson will take body as input & unmarshal it to res,
func UnmarshalJson(body string) (BPOD, error) {
	list := List{}
	res := BPOD{}

	err := json.Unmarshal([]byte(body), &list)
	if err != nil {
		return res, fmt.Errorf("UnmarshalJson failed\n%s", err.Error())
	}

	res = list.Photos[rand.Intn(len(list.Photos))]
	return res, nil
}

// GetJson takes reqInfo as input and returns the body and an error.
func GetJson(reqInfo map[string]string) (string, error) {
	// reqInfo is map[string]string and params is built from it,
	// currently it takes apiKey and the date from reqInfo to
	// build param. If any new key/value is added to reqInfo then
	// it must be addded here too, it won't be sent as param
	// directly.
	params := make(map[string]string)
	params["format"] = "js"
	params["n"] = "1"

	// if random is true then fetch 7 photos
	if reqInfo["random"] == "true" {
		params["n"] = "7"
	}

	body, err := request.GetRes(reqInfo["api"], params)
	return string(body), err
}

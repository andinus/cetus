package bpod

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"framagit.org/andinus/cetus/pkg/cetus"
)

type Res struct {
	StartDate     string `json:"startdate"`
	FullStartDate string `json:"fullstartdate"`
	EndDate       string `json:"enddate"`
	Url           string `json:"url"`
	UrlBase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Hsh           string `json:"hsh"`
}

type List struct {
	Photos []Res `json:"images"`
}

// UnmarshalJson will take body as input & unmarshal it to res
func UnmarshalJson(body string) (Res, error) {
	list := List{}
	res := Res{}

	err := json.Unmarshal([]byte(body), &list)
	if err != nil {
		return res, fmt.Errorf("UnmarshalJson failed\n%s", err.Error())
	}

	res = list.Photos[rand.Intn(len(list.Photos))]
	return res, nil
}

// GetJson returns json response received from the api
func GetJson(reqInfo map[string]string) (string, error) {
	params := make(map[string]string)
	params["format"] = "js"
	params["n"] = "1"

	// if random is true then fetch 7 photos
	if reqInfo["random"] == "true" {
		params["n"] = "7"

	}

	body, err := cetus.GetRes(reqInfo["api"], params)
	return string(body), err
}

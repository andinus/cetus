package bing

import (
	"time"

	"framagit.org/andinus/cetus/pkg/cetus"
)

// GetBpodJson returns json response received from the api
func GetBpodJson(reqInfo map[string]string, t time.Duration) (string, error) {
	params := make(map[string]string)
	params["format"] = "js"
	params["n"] = "1"

	// if random is true then fetch 7 photos
	if reqInfo["random"] == "true" {
		params["n"] = "7"

	}

	body, err := cetus.GetRes(reqInfo["api"], params, t)
	return string(body), err
}

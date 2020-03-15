package wallhaven

import (
	"time"

	"framagit.org/andinus/cetus/pkg/cetus"
)

// GetWhJson returns json response received from the api
func GetWhJson(reqInfo map[string]string, t time.Duration) (string, error) {
	params := make(map[string]string)
	params["apikey"] = reqInfo["apiKey"]
	if reqInfo["random"] == "true" {
		params["sorting"] = "random"

	}

	body, err := cetus.GetRes(reqInfo["api"], params, t)
	return string(body), err
}

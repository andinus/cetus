// Request manages all outgoing requests for cetus projects.
package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetRes takes api and params as input and returns the body and
// error.
func (c http.Client) GetRes(api string, params map[string]string) (string, error) {
	var body string

	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"request.go: failed to create request",
			err.Error())
		return body, err
	}

	// User-Agent should be passed with every request to make work
	// easier for the server handler. Include contact information
	// along with the project name so they could reach you if
	// required.
	req.Header.Set("User-Agent",
		"Andinus / Cetus - https://andinus.nand.sh/projects/cetus")

	// Params is a simple map[string]string which contains
	// parameters that needs to be passed along with the request.
	// There is no check involved here & it should be done before
	// passing params to this function.
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	res, err := c.Do(req)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"request.go: failed to get response",
			err.Error())
		return body, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("Unexpected response status code received: %d %s",
			res.StatusCode,
			http.StatusText(res.StatusCode))
		return body, err
	}

	// This will read everything to memory and is okay to use here
	// because the json response received will be small unlike in
	// download.go (package background) where it is an image.
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"request.go: failed to read body to out (var)",
			err.Error())
		return body, err
	}

	body = string(out)
	return body, err
}

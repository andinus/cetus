package background

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Download takes path and url as input and downloads the data to a
// file, returning an error if there is one.
func Download(file string, url string) error {
	o, err := os.Create(file)
	if err != nil {
		err = fmt.Errorf("%s%s\n%s",
			"download.go: failed to create file: ", file,
			err.Error())
		return err
	}
	defer o.Close()

	res, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("%s%s\n%s",
			"download.go: failed to get response from ", url,
			err.Error())
		return err
	}
	defer res.Body.Close()

	// Return an error on unexpected response code.
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unexpected Response: %s",
			res.Status)
		return err
	}

	// This will not copy everything to memory but will save to
	// disk as it progresses, ideal for big files or low memory
	// environments.
	_, err = io.Copy(o, res.Body)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"download.go: failed to copy body to file",
			err.Error())
	}
	return err
}

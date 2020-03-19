package background

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Download takes path and url as input and downloads the data to a
// file, returning an error if there is one
func Download(file string, url string) error {
	o, err := os.Create(file)
	if err != nil {
		return err
	}
	defer o.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected Response: %s", res.Status)
	}

	_, err = io.Copy(o, res.Body)
	if err != nil {
		return err
	}
	return nil
}

// Copyright (c) 2020, Andinus <andinus@inventati.org>

// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.

// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var (
	err     error
	timeout time.Duration
)

func main() {
	var (
		imgPath string

		apod       bool
		apodAPI    string
		apodAPIKey string
	)

	// Parse flags passed to program
	flag.StringVar(&imgPath, "img-path", "", "Image to set as wallpaper")

	flag.BoolVar(&apod, "apod", false, "Set Astronomy Picture of the Day as wallpaper")
	flag.StringVar(&apodAPI, "apod-api", "https://api.nasa.gov/planetary/apod", "APOD API URL")
	flag.StringVar(&apodAPIKey, "apod-api-key", "DEMO_KEY", "APOD API Key")

	flag.DurationVar(&timeout, "timeout", 16, "Timeout for http client")

	flag.Parse()

	if len(imgPath) > 0 {
		err = setWall(imgPath)
		errChk(err)
		return
	}

	if apod != false {
		apodI := make(map[string]string)
		apodI["api"] = apodAPI
		apodI["apiKey"] = apodAPIKey

		err = setWallFromAPOD(apodI)
		errChk(err)
		return
	}
}

// Calls feh to set the wallpaper
func setWall(imgPath string) error {
	feh, err := exec.LookPath("feh")
	if err != nil {
		fmt.Println("Error: feh is not in $PATH")
		return err
	}

	fmt.Printf("Path to set as Wallpaper: %s\n", imgPath)

	err = exec.Command(feh, "--bg-fill", imgPath).Run()
	return err
}

// Get url of Astronomy Picture of the Day & pass it to setWall()
func setWallFromAPOD(apodI map[string]string) error {
	type apodRes struct {
		Copyright      string `json:"copyright"`
		Date           string `json:"string"`
		Explanation    string `json:"explanation"`
		HDURL          string `json:"hdurl"`
		MediaType      string `json:"media_type"`
		ServiceVersion string `json:"service_version"`
		Title          string `json:"title"`
		URL            string `json:"url"`
	}

	apodNow := apodRes{}

	req, err := http.NewRequest(http.MethodGet, apodI["api"], nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("api_key", apodI["apiKey"])
	req.URL.RawQuery = q.Encode()

	res, err := getRes(req)
	if err != nil {
		fmt.Printf("Error: GET %s\n", apodI["api"])
		return err
	}
	defer res.Body.Close()

	apiBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(apiBody), &apodNow)
	if err != nil {
		return err
	}

	// Set Astronomy Picture of the Day as wallpaper
	err = setWall(apodNow.HDURL)
	return err
}

func getRes(req *http.Request) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * timeout,
	}
	res, err := client.Do(req)

	return res, err
}

func errChk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

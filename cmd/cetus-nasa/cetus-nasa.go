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
	"flag"
	"fmt"
	"log"
	"time"

	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/cetus"
	"framagit.org/andinus/cetus/pkg/nasa"
)

var (
	quiet     bool
	version   bool
	fetchOnly bool

	api         string
	apiKey      string
	date        string
	random      bool
	dateHelp    string
	dateDefault string
	timeout     time.Duration
)

func main() {
	parseFlags()

	if version {
		cetus.Version()
		return
	}

	// Convert timeout to seconds
	timeout = timeout * time.Second

	if random {
		date = nasa.RandDate()
	}

	var (
		picturePath string
		apodRes     nasa.APOD
		err         error
		apodInfo    map[string]string
	)

	apodInfo = make(map[string]string)
	apodInfo["api"] = api
	apodInfo["apiKey"] = apiKey
	apodInfo["date"] = date

	apodRes, err = nasa.APODPath(apodInfo, timeout)
	if err != nil {
		if len(apodRes.Msg) != 0 {
			log.Println("Message: ", apodRes.Msg)
		}
		log.Fatal(err)
	}

	if !quiet {
		printDetails(apodRes)
	}

	// if fetchOnly is true then don't set background
	if fetchOnly {
		return
	}
	picturePath = apodRes.HDURL

	// if media type is an image then set background
	if apodRes.MediaType == "image" {
		err = background.Set(picturePath)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parseFlags() {
	flag.BoolVar(&quiet, "quiet", false, "No output")
	flag.BoolVar(&version, "version", false, "Cetus version")
	flag.BoolVar(&fetchOnly, "fetch-only", false, "Don't set background, only fetch info")

	dateHelp = fmt.Sprintf("Choose a random date between 1995-06-16 & %s",
		time.Now().UTC().Format("2006-01-02"))
	flag.BoolVar(&random, "random", false, dateHelp)

	flag.StringVar(&api, "api", "https://api.nasa.gov/planetary/apod", "APOD API URL")
	flag.StringVar(&apiKey, "api-key", "DEMO_KEY", "api.nasa.gov key for expanded usage")

	dateDefault = time.Now().UTC().Format("2006-01-02")
	flag.StringVar(&date, "date", dateDefault, "Date of the APOD image to retrieve")

	flag.DurationVar(&timeout, "timeout", 32*time.Second, "Timeout for http client in seconds")
	flag.Parse()

}

func printDetails(apodRes nasa.APOD) {
	fmt.Printf("Title: %s\n\n", apodRes.Title)
	fmt.Printf("Copyright: %s\n", apodRes.Copyright)
	fmt.Printf("Date: %s\n\n", apodRes.Date)
	fmt.Printf("Media Type: %s\n", apodRes.MediaType)
	fmt.Printf("URL: %s\n\n", apodRes.HDURL)
	fmt.Printf("Explanation: %s\n", apodRes.Explanation)
}

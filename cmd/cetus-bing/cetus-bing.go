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
	"math/rand"
	"time"

	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/bing"
	"framagit.org/andinus/cetus/pkg/cetus"
)

var (
	quiet     bool
	version   bool
	fetchOnly bool
	pathOnly  bool

	api     string
	random  bool
	timeout time.Duration

	err       error
	bpodRes   bing.BPOD
	bpodPhoto bing.Photo
)

func main() {
	parseFlags()

	if version {
		cetus.Version()
		return
	}
	rand.Seed(time.Now().Unix())

	// Convert timeout to seconds
	timeout = timeout * time.Second

	// get response from api
	bpodRes, err = getBPODRes()
	if err != nil {
		log.Fatal(err)
	}

	// if random was set then bpodRes holds list of multiple
	// responses, choose a random response from the list
	var i int = rand.Intn(len(bpodRes.Photos))
	bpodPhoto = bpodRes.Photos[i]
	bpodPhoto.URL = fmt.Sprintf("%s%s", "https://www.bing.com", bpodPhoto.URL)
	printDetails(bpodPhoto)

	// if fetchOnly is true then don't set background
	if fetchOnly {
		return
	}

	// if media type is an image then set background
	err = background.Set(bpodPhoto.URL)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.BoolVar(&quiet, "quiet", false, "No output")
	flag.BoolVar(&version, "version", false, "Cetus version")
	flag.BoolVar(&fetchOnly, "fetch-only", false, "Don't set background, only fetch info")

	flag.BoolVar(&random, "random", false, "Choose a random image (from 7 images)")
	flag.BoolVar(&pathOnly, "path-only", false, "Print only path of the image")

	flag.StringVar(&api, "api", "https://www.bing.com/HPImageArchive.aspx", "BPOD API URL")

	flag.DurationVar(&timeout, "timeout", 32*time.Second, "Timeout for http client in seconds")
	flag.Parse()

}

func printDetails(bpodPhoto bing.Photo) {
	if quiet {
		return
	}
	if pathOnly {
		cetus.PrintPath(bpodPhoto.URL)
		return
	}
	fmt.Printf("Title: %s\n\n", bpodPhoto.Title)
	fmt.Printf("Copyright: %s\n", bpodPhoto.Copyright)
	fmt.Printf("Copyright Link: %s\n", bpodPhoto.CopyrightLink)
	fmt.Printf("Date: %s\n\n", bpodPhoto.StartDate)
	fmt.Printf("URL: %s\n\n", bpodPhoto.URL)
}

func getBPODRes() (bing.BPOD, error) {
	var bpodInfo map[string]string
	bpodInfo = make(map[string]string)
	bpodInfo["api"] = api
	if random {
		bpodInfo["random"] = "true"
	}
	bpodRes, err = bing.BPODPath(bpodInfo, timeout)

	return bpodRes, err
}

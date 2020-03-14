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
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

var (
	timeout time.Duration

	unsplashAPI string

	width  int
	height int
)

func main() {
	rand.Seed(time.Now().Unix())

	var (
		err error

		imgPath string
		wall    string
		src     string
		srcArr  []string = []string{
			"unsplash",
		}
	)

	// Parse flags passed to program
	flag.StringVar(&src, "src", "random", "Source for the image")
	flag.StringVar(&wall, "wall", "random", "Daily, Weekly or Random wallpaper")

	flag.IntVar(&width, "width", 1920, "Width of the image")
	flag.IntVar(&height, "height", 1080, "Height of the image")

	flag.StringVar(&unsplashAPI, "unsplash-api", "https://source.unsplash.com", "Unsplash Source API URL")
	flag.DurationVar(&timeout, "timeout", 16, "Timeout for http client")
	flag.Parse()

	if src == "random" {
		src = srcArr[rand.Intn(len(srcArr))]
	}

	// Check if the source is known
	if !contains(srcArr, src) {
		log.Fatal("Error: Unknown Source")
	}

	imgPath, err = parseSrcAndGetPath(src, wall)
	errChk(err)

	err = setWall(imgPath)
	errChk(err)
}

func contains(arr []string, str string) bool {
	for _, i := range arr {
		if i == str {
			return true
		}
	}
	return false
}

// Gets image path from src
func parseSrcAndGetPath(src string, wall string) (string, error) {
	var err error
	var imgPath string

	switch src {
	case "unsplash":
		fmt.Println("Unsplash Source")
		imgPath, err = getPathUnsplash(wall)
	}

	return imgPath, err
}

func getPathUnsplash(wall string) (string, error) {
	var err error
	var imgPath string

	switch wall {
	case "daily", "weekly":
		unsplashAPI = fmt.Sprintf("%s/%s",
			unsplashAPI, wall)
	case "random":
		unsplashAPI = fmt.Sprintf("%s/%sx%s",
			unsplashAPI, strconv.Itoa(width), strconv.Itoa(height))
	default:
		return "", fmt.Errorf("Error: Unknown wall")
	}

	req, err := http.NewRequest(http.MethodGet, unsplashAPI, nil)
	if err != nil {
		return "", err
	}

	res, err := getRes(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Unsplash Source API will redirect to the image
	imgPath = res.Request.URL.String()
	return imgPath, err
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

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"tildegit.org/andinus/cetus/background"
	"tildegit.org/andinus/cetus/bpod"
	"tildegit.org/andinus/cetus/cache"
	"tildegit.org/andinus/cetus/notification"
)

func execBPOD() {
	bpodApi := getEnv("BPOD_API", "https://www.bing.com/HPImageArchive.aspx")

	// reqInfo holds all the parameters that needs to be sent with
	// the request. GetJson() will pack apiKey & date in params
	// map before sending it to another function. Adding params
	// here will not change the behaviour of the function, changes
	// have to be made in GetJson() too.
	reqInfo = make(map[string]string)
	reqInfo["api"] = bpodApi

	if random {
		reqInfo["random"] = "true"
	}

	cacheDir := fmt.Sprintf("%s/%s", cache.GetDir(), "bpod")
	os.MkdirAll(cacheDir, os.ModePerm)

	body, err = bpod.GetJson(reqInfo)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"bpod.go: failed to get json response from api",
			err.Error())
		log.Fatal(err)
	}

	if dump {
		fmt.Println(body)
	}

	res, err := bpod.UnmarshalJson(body)
	if err != nil {
		log.Fatal(err)
	}

	// Correct format
	res.URL = fmt.Sprintf("%s%s", "https://www.bing.com", res.URL)
	dt, err := time.Parse("20060102", res.StartDate)
	if err != nil {
		log.Fatal(err)
	}
	res.StartDate = dt.Format("2006-01-02")

	// Send a desktop notification if notify flag was passed.
	if notify {
		n := notification.Notif{}
		n.Title = res.Title
		n.Message = fmt.Sprintf("%s\n\n%s",
			res.StartDate,
			res.Copyright)

		err = n.Notify()
		if err != nil {
			log.Println(err)
		}
	}

	if print {
		fmt.Printf("Title: %s\n\n", res.Title)
		fmt.Printf("Copyright: %s\n", res.Copyright)
		fmt.Printf("Copyright Link: %s\n", res.CopyrightLink)
		fmt.Printf("Date: %s\n\n", res.StartDate)
		fmt.Printf("URL: %s\n", res.URL)
	}

	// Proceed only if the command was set because if it was fetch
	// then it's already finished & should exit now.
	if os.Args[1] == "fetch" {
		os.Exit(0)
	}

	// Try to set background only if the media type is an image.
	// First it downloads the image to the cache directory and
	// then tries to set it with feh. If the download fails then
	// it exits with a non-zero exit code.
	imgFile := fmt.Sprintf("%s/%s:%s", cacheDir, res.StartDate, res.Title)

	// Check if the file is available locally, if it is then don't
	// download it again and set it from disk
	if _, err := os.Stat(imgFile); os.IsNotExist(err) {
		err = background.Download(imgFile, res.URL)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		if err != nil {
			log.Fatal(err)
		}
	}

	err = background.SetFromFile(imgFile)
	if err != nil {
		log.Fatal(err)
	}
}

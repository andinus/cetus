package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"tildegit.org/andinus/cetus/apod"
	"tildegit.org/andinus/cetus/background"
	"tildegit.org/andinus/cetus/cache"
	"tildegit.org/andinus/cetus/notification"
)

var (
	err     error
	body    string
	file    string
	reqInfo map[string]string
)

func execAPOD() {
	apodApi := getEnv("APOD_API", "https://api.nasa.gov/planetary/apod")
	apodKey := getEnv("APOD_KEY", "DEMO_KEY")

	// reqInfo holds all the parameters that needs to be sent with
	// the request. GetJson() will pack apiKey & date in params
	// map before sending it to another function. Adding params
	// here will not change the behaviour of the function, changes
	// have to be made in GetJson() too.
	reqInfo = make(map[string]string)
	reqInfo["api"] = apodApi
	reqInfo["apiKey"] = apodKey
	reqInfo["date"] = apodDate

	if random {
		reqInfo["date"] = apod.RandDate()
	}

	cacheDir := fmt.Sprintf("%s/%s", cache.GetDir(), "apod")
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check if the file is available locally, if it is then don't
	// download it again and get it from disk
	file = fmt.Sprintf("%s/%s.json", cacheDir, reqInfo["date"])

	if _, err := os.Stat(file); err == nil {
		data, err := ioutil.ReadFile(file)

		// Not being able to read from the cache file is a
		// small error and the program shouldn't exit but
		// should continue after printing the log so that the
		// user can investigate it later.
		if err != nil {
			err = fmt.Errorf("%s%s\n%s",
				"apod.go: failed to read file to data: ", file,
				err.Error())
			log.Println(err)
			dlAndCacheAPODBody()
		}
		body = string(data)

	} else if os.IsNotExist(err) {
		dlAndCacheAPODBody()

	} else {
		// If file existed then that is handled by the if
		// block, if it didn't exist then that is handled by
		// the else if block. If we reach here then that means
		// it's Schrödinger's file & something else went
		// wrong.
		log.Fatal(err)
	}

	if dump {
		fmt.Printf(body)
	}

	res := apod.APOD{}
	err = apod.UnmarshalJson(&res, body)
	if err != nil {
		log.Fatal(err)
	}

	// res.Msg will be returned when there is error on user input
	// or the api server.
	if len(res.Msg) != 0 {
		fmt.Printf("Message: %s", res.Msg)
		os.Exit(1)
	}

	// Send a desktop notification if notify flag was passed.
	if notify {
		n := notification.Notif{}
		n.Title = res.Title
		n.Message = fmt.Sprintf("%s\n\n%s",
			res.Date,
			res.Explanation)

		err = n.Notify()
		if err != nil {
			log.Println(err)
		}
	}

	if print {
		fmt.Printf("Title: %s\n\n", res.Title)
		if len(res.Copyright) != 0 {
			fmt.Printf("Copyright: %s\n", res.Copyright)
		}
		fmt.Printf("Date: %s\n\n", res.Date)
		fmt.Printf("Media Type: %s\n", res.MediaType)
		if res.MediaType == "image" {
			fmt.Printf("URL: %s\n\n", res.HDURL)
		} else {
			fmt.Printf("URL: %s\n\n", res.URL)
		}
		fmt.Printf("Explanation: %s\n", res.Explanation)
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
	if res.MediaType != "image" {
		os.Exit(0)
	}
	imgFile := fmt.Sprintf("%s/%s", cacheDir, res.Title)

	// Check if the file is available locally, if it is then don't
	// download it again and set it from disk.
	if _, err := os.Stat(imgFile); os.IsNotExist(err) {
		err = background.Download(imgFile, res.HDURL)
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

func dlAndCacheAPODBody() {
	body, err = apod.GetJson(reqInfo)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"apod.go: failed to get json response from api",
			err.Error())
		log.Fatal(err)
	}

	// Write body to the cache so that it can be read later.
	err = ioutil.WriteFile(file, []byte(body), 0644)

	// Not being able to write to the cache file is a small error
	// and the program shouldn't exit but should continue after
	// printing the log so that the user can investigate it later.
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"apod.go: failed to write body to file: ", file,
			err.Error())
		log.Println(err)
	}
}

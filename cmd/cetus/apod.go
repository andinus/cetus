package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"framagit.org/andinus/cetus/pkg/apod"
	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/indus/notification"
)

func execAPOD() {
	// reqInfo holds all the parameters that needs to be sent with
	// the request. GetJson() will pack apiKey & date in params
	// map before sending it to another function. Adding params
	// here will not change the behaviour of the function, changes
	// have to be made in GetJson() too.
	reqInfo := make(map[string]string)
	reqInfo["api"] = string(*apodAPI)
	reqInfo["apiKey"] = string(*apodKey)
	reqInfo["date"] = string(*apodDate)

	if *apodRand {
		reqInfo["date"] = apod.RandDate()
	}

	cacheDir := fmt.Sprintf("%s/%s", getCacheDir(), "apod")
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check if the file is available locally, if it is then don't
	// download it again and get it from disk
	var body string
	file := fmt.Sprintf("%s/%s", cacheDir, reqInfo["date"])
	if _, err := os.Stat(file); err == nil {
		data, err := ioutil.ReadFile(file)
		chkErr(err)
		body = string(data)
	} else if os.IsNotExist(err) {
		body, err = apod.GetJson(reqInfo)
		chkErr(err)

		// Write body to the cache so that it can be read
		// later
		err = ioutil.WriteFile(file, []byte(body), 0644)
		chkErr(err)
	} else {
		chkErr(err)
	}

	if *apodDump {
		fmt.Printf(body)
		os.Exit(0)
	}

	res := apod.Res{}
	err := apod.UnmarshalJson(&res, body)
	chkErr(err)

	// res.Msg will be returned when there is error on user input
	// or the api server.
	if len(res.Msg) != 0 {
		fmt.Printf("Message: %s", res.Msg)
		os.Exit(1)
	}

	// If path-only is passed then it will only print the path,
	// even if quiet is passed. If the user wants the program to
	// be quiet then path-only shouldn't be passed. If path-only
	// is not passed & quiet is also not passed then print the
	// response.
	//
	// Path is only printed when the media type is an image
	// because res.HDURL is empty on non image media type.
	if *apodPathOnly {
		fmt.Println(res.HDURL)
	} else if !*apodQuiet {
		apod.Print(res)
	}

	// Send a desktop notification if notify flag was passed
	if *apodNotify {
		n := notification.Notif{}
		n.Title = res.Title
		n.Message = fmt.Sprintf("%s\n\n%s",
			res.Date,
			res.Explanation)

		err = n.Notify()
		chkErr(err)
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
	imgCacheDir := fmt.Sprintf("%s/%s", cacheDir, "background")
	os.MkdirAll(imgCacheDir, os.ModePerm)
	imgFile := fmt.Sprintf("%s/%s", imgCacheDir, reqInfo["date"])

	// Check if the file is available locally, if it is
	// then don't download it again and set it from disk
	if _, err := os.Stat(imgFile); os.IsNotExist(err) {
		err = background.Download(imgFile, res.HDURL)
		chkErr(err)
	} else {
		chkErr(err)
	}

	err = background.Set(imgFile)
	chkErr(err)
}

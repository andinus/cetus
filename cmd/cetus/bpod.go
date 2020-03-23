package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/bpod"
	"framagit.org/andinus/indus/notification"
)

func execBPOD() {
	// reqInfo holds all the parameters that needs to be sent with
	// the request. GetJson() will pack apiKey & date in params
	// map before sending it to another function. Adding params
	// here will not change the behaviour of the function, changes
	// have to be made in GetJson() too.
	reqInfo := make(map[string]string)
	reqInfo["api"] = string(*bpodAPI)

	if *bpodRand {
		reqInfo["random"] = "true"
	}

	cacheDir := fmt.Sprintf("%s/%s", getCacheDir(), "bpod")
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check if the file is available locally, if it is then don't
	// download it again and get it from disk.
	//
	// We don't know the bpod date because that will be there in
	// response & we can't read the response without requesting
	// it. So this will assume the bpod date to be today's date if
	// *bpodRand is not set true. If *bpodRand is set true then we
	// can't assume the date. Also this way too it can cause error
	// if our assumed date doesn't matches date at the server.
	var body string
	var file string
	var err error

	if !*bpodRand {
		// If not *bpodRand and the file exists then read from
		// disk, if the file doesn't exist then get it and
		// save it to disk.
		file = fmt.Sprintf("%s/%s", cacheDir, time.Now().UTC().Format("2006-01-02"))
		if _, err := os.Stat(file); err == nil {
			data, err := ioutil.ReadFile(file)
			chkErr(err)
			body = string(data)
		} else if os.IsNotExist(err) {
			body, err = bpod.GetJson(reqInfo)
			chkErr(err)

			// Write body to the cache so that it can be
			// read later
			err = ioutil.WriteFile(file, []byte(body), 0644)
			chkErr(err)
		} else {
			chkErr(err)
		}
	} else {
		// If *bpodRand then get the file and save it to disk
		// after unmarshal because we don't know the file name
		// yet
		body, err = bpod.GetJson(reqInfo)
		chkErr(err)
	}

	// Unmarshal before dump because otherwise if we come across
	// the date for the first time then it would just dump and
	// exit without saving it to cache. This way we first save it
	// to cache if *bpodRand is true.
	res, err := bpod.UnmarshalJson(body)
	chkErr(err)

	// Correct format
	res.Url = fmt.Sprintf("%s%s", "https://www.bing.com", res.Url)
	dt, err := time.Parse("20060102", res.StartDate)
	chkErr(err)
	res.StartDate = dt.Format("2006-01-02")

	file = fmt.Sprintf("%s/%s", cacheDir, res.StartDate)
	if *bpodRand {
		// Write body to the cache so that it can be read
		// later
		err = ioutil.WriteFile(file, []byte(body), 0644)
		chkErr(err)
	}

	if *bpodDump {
		fmt.Printf(body)
		os.Exit(0)
	}

	// If path-only is passed then it will only print the path,
	// even if quiet is passed. If the user wants the program to
	// be quiet then path-only shouldn't be passed. If path-only
	// is not passed & quiet is also not passed then print the
	// response.
	//
	// Path is only printed when the media type is an image
	// because res.HDURL is empty on non image media type.
	if *bpodPathOnly {
		fmt.Println(res.Url)
	} else if !*bpodQuiet {
		bpod.Print(res)
	}

	// Send a desktop notification if notify flag was passed
	if *bpodNotify {
		n := notification.Notif{}
		n.Title = res.Title
		n.Message = fmt.Sprintf("%s\n\n%s",
			res.StartDate,
			res.Copyright)

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
	imgCacheDir := fmt.Sprintf("%s/%s", cacheDir, "background")
	os.MkdirAll(imgCacheDir, os.ModePerm)
	imgFile := fmt.Sprintf("%s/%s", imgCacheDir, res.StartDate)

	// Check if the file is available locally, if it is then don't
	// download it again and set it from disk
	if _, err := os.Stat(imgFile); os.IsNotExist(err) {
		err = background.Download(imgFile, res.Url)
		chkErr(err)
	} else {
		chkErr(err)
	}

	err = background.Set(imgFile)
	chkErr(err)
}

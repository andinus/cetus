package main

import (
	"fmt"
	"io/ioutil"
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

	// Check if the file is available locally, if it is then don't
	// download it again and get it from disk
	//
	// We don't know the bpod date because that will be there in
	// response & we can't read the response without requesting
	// it. So this will assume the bpod date to be today's date if
	// *bpodRand is not set true. If *bpodRand is set true then we
	// can't assume the date. Also this way too it can cause error
	// if our assumed date doesn't matches date at the server.

	if !random {
		// If not random and the file exists then read from
		// disk, if the file doesn't exist then get it and
		// save it to disk.
		file = fmt.Sprintf("%s/%s.json", cacheDir, time.Now().UTC().Format("2006-01-02"))
		if _, err := os.Stat(file); err == nil {
			data, err := ioutil.ReadFile(file)
			// Not being able to read from the cache file
			// is a small error and the program shouldn't
			// exit but should continue after printing the
			// log so that the user can investigate it
			// later.
			if err != nil {
				err = fmt.Errorf("%s%s\n%s",
					"bpod.go: failed to read file to data: ", file,
					err.Error())
				log.Println(err)
				dlAndCacheBPODBody()
			}
			body = string(data)

		} else if os.IsNotExist(err) {
			dlAndCacheBPODBody()

		} else {
			// If file existed then that is handled by the
			// if block, if it didn't exist then that is
			// handled by the else if block. If we reach
			// here then that means it's Schrödinger's
			// file & something else went wrong.
			log.Fatal(err)
		}
	} else {
		// If random then get the file and save it to disk
		// after unmarshal because we don't know the file name
		// yet
		body, err = bpod.GetJson(reqInfo)
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"bpod.go: failed to get json response from api",
				err.Error())
			log.Fatal(err)
		}

	}
	// Unmarshal before dump because otherwise if we come across
	// the date for the first time then it would just dump and
	// exit without saving it to cache. This way we first save it
	// to cache if *bpodRand is true.
	res, err := bpod.UnmarshalJson(body)
	if err != nil {
		log.Fatal(err)
	}

	// res and body are out of sync if random was passed because
	// body has all 7 entries but res has only one because
	// UnmarshalJson returns only one entry selected randomly.
	// Look at comment in UnmarshalJson and MarshalJson to
	// understand why this is required.
	if random {
		body, err = bpod.MarshalJson(res)

		// If an error was caused then that means body and res
		// is still out of sync which is not a big issue and
		// program shouldn't stop because of this, instead we
		// set random=false so that this body doesn't get
		// saved to cache and also add warning if user passed
		// the dump flag because dump is dumping body which
		// has all 7 entries.
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"bpod.go: failed to marshal res to body, both out of sync",
				err.Error())
			log.Println(err)

			log.Println("bpod.go: not saving this to cache")
			log.Println("bpod.go: dump will contain incorrect information")

			// This will prevent the program from saving
			// this to cache.
			random = false
		}
	}

	// Correct format
	res.URL = fmt.Sprintf("%s%s", "https://www.bing.com", res.URL)
	dt, err := time.Parse("20060102", res.StartDate)
	if err != nil {
		log.Fatal(err)
	}
	res.StartDate = dt.Format("2006-01-02")

	file = fmt.Sprintf("%s/%s.json", cacheDir, res.StartDate)
	if random {

		// Write body to the cache so that it can be read
		// later.
		err = ioutil.WriteFile(file, []byte(body), 0644)

		// Not being able to write to the cache file is a
		// small error and the program shouldn't exit but
		// should continue after printing the log so that the
		// user can investigate it later.
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"bpod.go: failed to write body to file: ", file,
				err.Error())
			log.Println(err)
		}
	}

	if dump {
		fmt.Printf(body)
	}

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
	imgFile := fmt.Sprintf("%s/%s", cacheDir, res.Title)

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

func dlAndCacheBPODBody() {
	body, err = bpod.GetJson(reqInfo)
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"bpod.go: failed to get json response from api",
			err.Error())
		log.Fatal(err)
	}

	// Write body to the cache so that it can be read later. This
	// will mess up the cache because the cache format will be
	// inconsistent, res is not the same as body and they should
	// be in sync before saving body to cache but here we don't
	// yet have res. This issue arises when the user passes random
	// flag, we are saving to cache the raw body returned but when
	// random is passed the format of body is changed to only
	// values unmarshalled in res because body is rebuilt from res
	// so the cache will have different format of body. Disabling
	// this cache for now seems like a good option, later we can
	// figure out how to make this body and body when random is
	// passed of same format.
	// err = ioutil.WriteFile(file, []byte(body), 0644)

	// Not being able to write to the cache file is a small error
	// and the program shouldn't exit but should continue after
	// printing the log so that the user can investigate it later.
	// if err != nil {
	// 	err = fmt.Errorf("%s\n%s",
	// 		"bpod.go: failed to write body to file: ", file,
	// 		err.Error())
	// 	log.Println(err)
	// }
}

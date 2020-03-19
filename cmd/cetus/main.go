package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"framagit.org/andinus/cetus/pkg/apod"
	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/bpod"
)

var (
	apodAPI      *string
	apodKey      *string
	apodDate     *string
	apodRand     *bool
	apodPathOnly *bool
	apodQuiet    *bool
	apodDump     *bool

	bpodAPI      *string
	bpodRand     *bool
	bpodPathOnly *bool
	bpodQuiet    *bool
	bpodDump     *bool
)

func main() {
	// Early Check: If command was not passed then print usage and
	// exit. Later command & service both are checked, this check
	// if for version command. If not checked then running cetus
	// without any args will fail because os.Args[1] will panic
	// the program & produce runtime error.
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	version := "v0.5.1"

	if os.Args[1] == "version" {
		fmt.Printf("Cetus %s\n", version)
		os.Exit(0)
	}

	// If command & service was not passed then print usage and
	// exit.
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())

	apodCmd := flag.NewFlagSet("apod", flag.ExitOnError)
	defDate := time.Now().UTC().Format("2006-01-02")

	// Flags to parse for apod service.
	apodAPI = apodCmd.String("api", "https://api.nasa.gov/planetary/apod", "APOD API link")
	apodKey = apodCmd.String("api-key", "DEMO_KEY", "NASA API Key for expanded usage")
	apodDate = apodCmd.String("date", defDate, "Date of NASA APOD to retrieve")
	apodRand = apodCmd.Bool("random", false, "Choose a date random starting from 1995-06-16")
	apodPathOnly = apodCmd.Bool("path-only", false, "Print only the path")
	apodQuiet = apodCmd.Bool("quiet", false, "Stay quiet")
	apodDump = apodCmd.Bool("dump", false, "Dump received response")

	bpodCmd := flag.NewFlagSet("bpod", flag.ExitOnError)

	// Flags to parse for bpod service.
	bpodAPI = bpodCmd.String("api", "https://www.bing.com/HPImageArchive.aspx", "BPOD API")
	bpodRand = bpodCmd.Bool("random", false, "Choose a random image from last week's BPOD")
	bpodPathOnly = bpodCmd.Bool("path-only", false, "Print only the path")
	bpodQuiet = bpodCmd.Bool("quiet", false, "Stay quiet")
	bpodDump = bpodCmd.Bool("dump", false, "Dump received response")

	// Switching on commands will cause more repetition than
	// switching on service. If we switch on commands then switch
	// on service will have to be replicated on every command
	// switch. Reverse is also true, this way we will repeat
	// command switch in every service but we do so in a better
	// way.
	//
	// However we check if the correct command was passed. version
	// command is not included because it has been dealt with
	// earlier in the program & the program should've exited after
	// that, if it reaches here then it's an error.
	switch os.Args[1] {
	case "set", "fetch":
	default:
		fmt.Printf("Invalid command: %q\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	switch os.Args[2] {
	case "apod", "nasa":
		apodCmd.Parse(os.Args[3:])
		if apodCmd.Parsed() {
			execAPOD()
		}
	case "bpod", "bing":
		bpodCmd.Parse(os.Args[3:])
		if bpodCmd.Parsed() {
			execBPOD()
		}
	default:
		fmt.Printf("Invalid service: %q\n", os.Args[2])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: cetus <command> <service> [<flags>]\n")
	fmt.Println("Commands: ")
	fmt.Println(" set     Set the latest image as background")
	fmt.Println(" fetch   Fetch the latest image information")
	fmt.Println(" version Print version")
	fmt.Println("\nServices: ")
	fmt.Println(" apod   NASA Astronomy Picture of the Day")
	fmt.Println(" bpod   Bing Photo of the Day")
}

// Check whether user has set CETUS_CACHE_DIR, if not then use the
// XDG_CACHE_HOME. If XDG_CACHE_HOME is not set then $HOME/.config
// should be used, according to XDG Base Directory Specification
func getCacheDir() string {
	cacheDir := os.Getenv("CETUS_CACHE_DIR")
	if len(cacheDir) == 0 {
		cacheDir = os.Getenv("XDG_CACHE_HOME")
	}
	if len(cacheDir) == 0 {
		cacheDir = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"),
			".cache", "cetus")
	}
	return cacheDir
}

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

func chkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

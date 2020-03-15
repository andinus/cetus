package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/cetus"
	"framagit.org/andinus/cetus/pkg/nasa"
)

type apod struct {
	copyright      string `json:"copyright"`
	date           string `json:"date"`
	explanation    string `json:"explanation"`
	hdURL          string `json:"hdurl"`
	mediaType      string `json:"media_type"`
	serviceVersion string `json:"service_version"`
	title          string `json:"title"`
	url            string `json:"url"`

	code int    `json:"code"`
	msg  string `json:"msg"`
}

var (
	t         time.Duration
	api       string
	date      string
	dump      bool
	quiet     bool
	random    bool
	apiKey    string
	version   bool
	fetchOnly bool
	pathOnly  bool
)

func main() {
	parseFlags()

	if version {
		cetus.Version()
		return
	}

	// Convert timeout to seconds
	t = t * time.Second

	if random {
		date = nasa.RandDate()
	}

	body, err := apodBody()
	if dump {
		fmt.Println(body)
		return
	}

	apod := apod{}
	err = json.Unmarshal([]byte(body), &apod)
	cetus.ErrChk("body unmarshal failed", err)
	if len(apod.msg) != 0 {
		log.Println("Message: ", apod.msg)
	}

	printDetails(apod)

	// if fetchOnly is true then don't set background
	if fetchOnly {
		return
	}

	// if media type is an image then set background
	if apod.mediaType == "image" {
		err = background.Set(apod.hdURL)
		cetus.ErrChk("setting background failed", err)
	}

}

func parseFlags() {
	flag.BoolVar(&quiet, "quiet", false, "No output")
	flag.BoolVar(&version, "version", false, "Cetus version")
	flag.BoolVar(&fetchOnly, "fetch-only", false, "Don't set background, only fetch info")
	flag.BoolVar(&dump, "dump", false, "Only dump received response")
	dateHelp := fmt.Sprintf("Choose a random date between 1995-06-16 & %s",
		time.Now().UTC().Format("2006-01-02"))
	flag.BoolVar(&random, "random", false, dateHelp)
	flag.BoolVar(&pathOnly, "path-only", false, "Print only path of the image")

	flag.StringVar(&api, "api", "https://api.nasa.gov/planetary/apod", "APOD API URL")
	flag.StringVar(&apiKey, "api-key", "DEMO_KEY", "api.nasa.gov key for expanded usage")

	dateDefault := time.Now().UTC().Format("2006-01-02")
	flag.StringVar(&date, "date", dateDefault, "Date of the APOD image to retrieve")

	flag.DurationVar(&t, "timeout", 32*time.Second, "Timeout for http client in seconds")
	flag.Parse()

}

func printDetails(apod apod) {
	if quiet {
		return
	}
	if pathOnly {
		cetus.PrintPath(apod.hdURL)
		return
	}
	fmt.Printf("Title: %s\n\n", apod.title)
	fmt.Printf("Copyright: %s\n", apod.copyright)
	fmt.Printf("Date: %s\n\n", apod.date)
	fmt.Printf("Media Type: %s\n", apod.mediaType)
	if apod.mediaType == "image" {
		fmt.Printf("URL: %s\n\n", apod.hdURL)
	} else {
		fmt.Printf("URL: %s\n\n", apod.url)
	}
	fmt.Printf("Explanation: %s\n", apod.explanation)
}

func apodBody() (string, error) {
	reqInfo := make(map[string]string)
	reqInfo["api"] = api
	reqInfo["apiKey"] = apiKey
	reqInfo["date"] = date

	body, err := nasa.GetApodJson(reqInfo, t)
	return body, err
}

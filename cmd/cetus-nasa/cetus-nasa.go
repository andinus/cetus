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
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HdURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	Url            string `json:"url"`

	Code int    `json:"code"`
	Msg  string `json:"msg"`
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
	if len(apod.Msg) != 0 {
		log.Println("Message: ", apod.Msg)
	}

	printDetails(apod)

	// if fetchOnly is true then don't set background
	if fetchOnly {
		return
	}

	// if media type is an image then set background
	if apod.MediaType == "image" {
		err = background.Set(apod.HdURL)
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
		cetus.PrintPath(apod.HdURL)
		return
	}
	fmt.Printf("Title: %s\n\n", apod.Title)
	fmt.Printf("Copyright: %s\n", apod.Copyright)
	fmt.Printf("Date: %s\n\n", apod.Date)
	fmt.Printf("Media Type: %s\n", apod.MediaType)
	if apod.MediaType == "image" {
		fmt.Printf("URL: %s\n\n", apod.HdURL)
	} else {
		fmt.Printf("URL: %s\n\n", apod.Url)
	}
	fmt.Printf("Explanation: %s\n", apod.Explanation)
}

func apodBody() (string, error) {
	reqInfo := make(map[string]string)
	reqInfo["api"] = api
	reqInfo["apiKey"] = apiKey
	reqInfo["date"] = date

	body, err := nasa.GetApodJson(reqInfo, t)
	return body, err
}

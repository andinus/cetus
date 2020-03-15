package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/bing"
	"framagit.org/andinus/cetus/pkg/cetus"
)

type photo struct {
	startDate     string `json:"startdate"`
	fullStartDate string `json:"fullstartdate"`
	endDate       string `json:"enddate"`
	url           string `json:"url"`
	urlBase       string `json:"urlbase"`
	copyright     string `json:"copyright"`
	copyrightLink string `json:"copyrightlink"`
	title         string `json:"title"`
	hsh           string `json:"hsh"`
}

type bpod struct {
	photos []photo `json:"images"`
}

var (
	t         time.Duration
	api       string
	dump      bool
	quiet     bool
	random    bool
	version   bool
	pathOnly  bool
	fetchOnly bool
)

func main() {
	parseFlags()

	if version {
		cetus.Version()
		return
	}
	rand.Seed(time.Now().Unix())

	// Convert timeout to seconds
	t = t * time.Second

	body, err := bpodBody()
	if dump {
		fmt.Println(body)
		return
	}

	bpod := bpod{}
	err = json.Unmarshal([]byte(body), &bpod)
	errChk("body unmarshal failed", err)

	// if random was set then bpodRes holds list of multiple
	// responses, choose a random response from the list
	var i int = rand.Intn(len(bpod.photos))
	bpodPhoto := bpod.photos[i]

	// correct image path
	bpodPhoto.url = fmt.Sprintf("%s%s", "https://www.bing.com", bpodPhoto.url)

	// correct date format
	dt, err := time.Parse("20060102", bpodPhoto.startDate)
	errChk("bpodPhoto.startDate parse failed", err)
	bpodPhoto.startDate = dt.Format("2006-01-02")

	printDetails(bpodPhoto)

	// if fetchOnly is true then don't set background
	if fetchOnly {
		return
	}

	err = background.Set(bpodPhoto.url)
	errChk("setting background failed", err)
}

func parseFlags() {
	flag.BoolVar(&quiet, "quiet", false, "No output")
	flag.BoolVar(&version, "version", false, "Cetus version")
	flag.BoolVar(&fetchOnly, "fetch-only", false, "Don't set background, only fetch info")
	flag.BoolVar(&dump, "dump", false, "Only dump received response")
	flag.BoolVar(&random, "random", false, "Choose a random image (from 7 images)")
	flag.BoolVar(&pathOnly, "path-only", false, "Print only path of the image")

	flag.StringVar(&api, "api", "https://www.bing.com/HPImageArchive.aspx", "BPOD API URL")

	flag.DurationVar(&t, "timeout", 32*time.Second, "Timeout for http client in seconds")
	flag.Parse()

}

func printDetails(bpodPhoto photo) {
	if quiet {
		return
	}
	if pathOnly {
		cetus.PrintPath(bpodPhoto.url)
		return
	}
	fmt.Printf("Title: %s\n\n", bpodPhoto.title)
	fmt.Printf("Copyright: %s\n", bpodPhoto.copyright)
	fmt.Printf("Copyright Link: %s\n", bpodPhoto.copyrightLink)
	fmt.Printf("Date: %s\n\n", bpodPhoto.startDate)
	fmt.Printf("URL: %s\n", bpodPhoto.url)
}

func bpodBody() (string, error) {
	reqInfo := make(map[string]string)
	reqInfo["api"] = api
	if random {
		reqInfo["random"] = "true"
	}

	body, err := bing.GetBpodJson(reqInfo, t)
	return body, err
}

func errChk(ctx string, err error) {
	log.Println(ctx)
	log.Fatal(err)
}

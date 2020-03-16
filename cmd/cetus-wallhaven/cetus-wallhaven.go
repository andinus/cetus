package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/cetus"
	"framagit.org/andinus/cetus/pkg/wallhaven"
)

type wh struct {
	Photos []photo `json:"data"`
	// MetaI  []meta  `json:"meta"`
}

// type meta struct {
// 	CurrentPage int    `json:"current_page"`
// 	LastPage    int    `json:"last_page"`
// 	PerPage     int    `json:"per_page"`
// 	Total       int    `json:"total"`
// 	Query       string `json:"query"`
// 	Seed        string `json:"seed"`
// }

type photo struct {
	Id         string   `json:"id"`
	Url        string   `json:"url"`
	ShortUrl   string   `json:"short_url"`
	Views      int      `json:"views"`
	Favorites  int      `json:"favorites"`
	Source     string   `json:"source"`
	Purity     string   `json:"purity"`
	Category   string   `json:"category"`
	DimensionX int      `json:"dimension_x"`
	DimensionY int      `json:"dimension_y"`
	Resolution string   `json:"resolution"`
	Ratio      string   `json:"ratio"`
	FileSize   int      `json:"file_size"`
	FileType   string   `json:"file_type"`
	CreatedAt  string   `json:"created_at"`
	Path       string   `json:"path"`
	Colors     []string `json:"colors"`
	// Thumbs     []thumb  `json:"thumbs"`
}

// type thumb struct {
// 	Small    string `json:"small"`
// 	Original string `json:"original"`
// 	Large    string `json:"large"`
// }

var (
	t         time.Duration
	api       string
	apiKey    string
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

	body, err := whBody()
	if dump {
		fmt.Println(body)
		return
	}

	wh := wh{}
	err = json.Unmarshal([]byte(body), &wh)
	cetus.ErrChk("body unmarshal failed", err)

	// if random was set then wh holds list of multiple responses,
	// choose a random response from the list
	i := rand.Intn(len(wh.Photos))
	whPhoto := wh.Photos[i]

	printDetails(whPhoto)

	// if fetchOnly is true then don't set background
	if fetchOnly {
		return
	}

	err = background.Set(whPhoto.Path)
	cetus.ErrChk("setting background failed", err)
}

func parseFlags() {
	flag.BoolVar(&quiet, "quiet", false, "No output")
	flag.BoolVar(&version, "version", false, "Cetus version")
	flag.BoolVar(&fetchOnly, "fetch-only", false, "Don't set background, only fetch info")
	flag.BoolVar(&dump, "dump", false, "Only dump received response")
	flag.BoolVar(&random, "random", true, "Choose a random image")
	flag.BoolVar(&pathOnly, "path-only", false, "Print only path of the image")

	flag.StringVar(&api, "api", "https://wallhaven.cc/api/v1/search", "Wallhaven Search API URL")
	flag.StringVar(&apiKey, "api-key", "", "Wallhaven API Key")

	flag.DurationVar(&t, "timeout", 32*time.Second, "Timeout for http client in seconds")
	flag.Parse()

}

func printDetails(whPhoto photo) {
	if quiet {
		return
	}
	if pathOnly {
		cetus.PrintPath(whPhoto.Path)
		return
	}
	fmt.Printf("Id: %s\n", whPhoto.Id)
	fmt.Printf("URL: %s\n", whPhoto.Url)
	fmt.Printf("Short URL: %s\n", whPhoto.ShortUrl)
	fmt.Printf("Source: %s\n", whPhoto.Source)
	fmt.Printf("Date: %s\n\n", whPhoto.CreatedAt)
	fmt.Printf("Resolution: %s\n", whPhoto.Resolution)
	fmt.Printf("Ratio: %s\n", whPhoto.Ratio)
	fmt.Printf("Views: %d\n", whPhoto.Views)
	fmt.Printf("Favorites: %d\n", whPhoto.Favorites)
	fmt.Printf("File Size: %d KiB\n", whPhoto.FileSize/1024)
	fmt.Printf("Category: %s\n", whPhoto.Category)
}

func whBody() (string, error) {
	reqInfo := make(map[string]string)
	reqInfo["api"] = api
	if random {
		reqInfo["random"] = "true"
	}

	body, err := wallhaven.GetWhJson(reqInfo, t)
	return body, err
}

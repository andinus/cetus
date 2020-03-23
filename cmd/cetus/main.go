package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
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

	apodNotify *bool
	bpodNotify *bool
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
	bpodCmd := flag.NewFlagSet("bpod", flag.ExitOnError)

	defDate := time.Now().UTC().Format("2006-01-02")

	// Flags to parse for apod service.
	apodAPI = apodCmd.String("api", "https://api.nasa.gov/planetary/apod", "APOD API link")
	apodKey = apodCmd.String("api-key", "DEMO_KEY", "NASA API Key for expanded usage")
	apodDate = apodCmd.String("date", defDate, "Date of NASA APOD to retrieve")
	apodRand = apodCmd.Bool("random", false, "Choose a date random starting from 1995-06-16")
	apodPathOnly = apodCmd.Bool("path-only", false, "Print only the path")
	apodQuiet = apodCmd.Bool("quiet", false, "Stay quiet")
	apodDump = apodCmd.Bool("dump", false, "Dump received response")
	apodNotify = apodCmd.Bool("notify", false, "Send a desktop notification with background information")

	// Flags to parse for bpod service.
	bpodAPI = bpodCmd.String("api", "https://www.bing.com/HPImageArchive.aspx", "BPOD API")
	bpodRand = bpodCmd.Bool("random", false, "Choose a random image from last week's BPOD")
	bpodPathOnly = bpodCmd.Bool("path-only", false, "Print only the path")
	bpodQuiet = bpodCmd.Bool("quiet", false, "Stay quiet")
	bpodDump = bpodCmd.Bool("dump", false, "Dump received response")
	bpodNotify = bpodCmd.Bool("notify", false, "Send a desktop notification with background information")

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

func chkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

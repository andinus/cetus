package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"framagit.org/andinus/cetus/pkg/apod"
	"framagit.org/andinus/cetus/pkg/background"
	"framagit.org/andinus/cetus/pkg/bpod"
	"framagit.org/andinus/cetus/pkg/cetus"
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

	if os.Args[1] == "version" {
		cetus.Version()
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
	apodAPI := apodCmd.String("api", "https://api.nasa.gov/planetary/apod", "APOD API link")
	apodKey := apodCmd.String("api-key", "DEMO_KEY", "NASA API Key for expanded usage")
	apodDate := apodCmd.String("date", defDate, "Date of NASA APOD to retrieve")
	apodRand := apodCmd.Bool("random", false, "Choose a date random starting from 1995-06-16")
	apodPathOnly := apodCmd.Bool("path-only", false, "Print only the path")
	apodQuiet := apodCmd.Bool("quiet", false, "Stay quiet")
	apodDump := apodCmd.Bool("dump", false, "Dump received response")

	bpodCmd := flag.NewFlagSet("bpod", flag.ExitOnError)

	// Flags to parse for bpod service.
	bpodAPI := bpodCmd.String("api", "https://www.bing.com/HPImageArchive.aspx", "BPOD API")
	bpodRand := bpodCmd.Bool("random", false, "Choose a random image from last week's BPOD")
	bpodPathOnly := bpodCmd.Bool("path-only", false, "Print only the path")
	bpodQuiet := bpodCmd.Bool("quiet", false, "Stay quiet")
	bpodDump := bpodCmd.Bool("dump", false, "Dump received response")

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
	case "bpod", "bing":
		bpodCmd.Parse(os.Args[3:])
	default:
		fmt.Printf("Invalid service: %q\n", os.Args[2])
		printUsage()
		os.Exit(1)
	}

	if apodCmd.Parsed() {
		// reqInfo holds all the parameters that needs to be
		// sent with the request. GetJson() will pack apiKey &
		// date in params map before sending it to another
		// function. Adding params here will not change the
		// behaviour of the function, changes have to be made
		// in GetJson() too.
		reqInfo := make(map[string]string)
		reqInfo["api"] = string(*apodAPI)
		reqInfo["apiKey"] = string(*apodKey)
		reqInfo["date"] = string(*apodDate)

		if *apodRand {
			reqInfo["apiKey"] = apod.RandDate()
		}

		body, err := apod.GetJson(reqInfo)
		chkErr(err)

		if *apodDump {
			fmt.Printf(body)
			os.Exit(0)
		}

		res := apod.Res{}
		err = apod.UnmarshalJson(&res, body)
		chkErr(err)

		// res.Msg will be returned when there is error on
		// user input or the api server.
		if len(res.Msg) != 0 {
			fmt.Printf("Message: %s", res.Msg)
			os.Exit(1)
		}

		// If path-only is passed then it will only print the
		// path, even if quiet is passed. If the user wants
		// the program to be quiet then path-only shouldn't be
		// passed. If path-only is not passed & quiet is also
		// not passed then print the response.
		//
		// Path is only printed when the media type is an
		// image because res.HDURL is empty on non image media
		// type.
		if *apodPathOnly {
			fmt.Println(res.HDURL)
		} else if !*apodQuiet {
			apod.Print(res)
		}

		// Proceed only if the command was set because if it
		// was fetch then it's already finished & should exit
		// now.
		if os.Args[1] == "fetch" {
			os.Exit(0)
		}

		// Try to set background only if the media type is an
		// image.
		if res.MediaType == "image" {
			err = background.Set(res.HDURL)
			chkErr(err)
		}
	}

	if bpodCmd.Parsed() {
		// reqInfo here works similar to apodCmd block's
		// reqInfo, refer to explanation there.
		reqInfo := make(map[string]string)
		reqInfo["api"] = string(*bpodAPI)

		if *bpodRand {
			reqInfo["random"] = "true"
		}

		body, err := bpod.GetJson(reqInfo)
		chkErr(err)

		if *bpodDump {
			fmt.Printf(body)
			os.Exit(0)
		}

		res, err := bpod.UnmarshalJson(body)
		chkErr(err)

		// Correct format
		res.Url = fmt.Sprintf("%s%s", "https://www.bing.com", res.Url)
		dt, err := time.Parse("20060102", res.StartDate)
		chkErr(err)
		res.StartDate = dt.Format("2006-01-02")

		// path-only here works similar to apodCmd block's
		// path-only, refer to explanation there.
		if *bpodPathOnly {
			fmt.Println(res.Url)
		} else if !*bpodQuiet {
			bpod.Print(res)
		}

		// Proceed only if the command was set because if it
		// was fetch then it's already finished & should exit
		// now.
		if os.Args[1] == "fetch" {
			os.Exit(0)
		}

		err = background.Set(res.Url)
		chkErr(err)
	}
}

func printUsage() {
	fmt.Println("Usage: cetus <command> <service> [<args>]\n")
	fmt.Println("Commands: ")
	fmt.Println(" set   Set the latest image as background")
	fmt.Println(" fetch Fetch the latest image information")
	fmt.Println(" version Print version")
	fmt.Println("Services: ")
	fmt.Println(" apod   NASA Astronomy Picture of the Day")
	fmt.Println(" bpod   Bing Photo of the Day")
}

func chkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

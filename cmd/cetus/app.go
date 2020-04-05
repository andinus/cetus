package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	version string = "v0.6.7"
	dump    bool
	random  bool
	notify  bool
	print   bool

	err     error
	body    string
	file    string
	reqInfo map[string]string

	apodDate string
)

func app() {
	// Early Check: If command was not passed then print usage and
	// exit. Later command & service both are checked, this check
	// is for version command. If not checked then running cetus
	// without any args will fail because os.Args[1] will panic
	// the program & produce runtime error.
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}

	parseArgs()
}

// parseArgs will be parsing the arguments, it will verify if they are
// correct. Flag values are also set by parseArgs.
func parseArgs() {
	// Running just `cetus` would've paniced the program if length
	// of os.Args was not checked beforehand because there would
	// be no os.Args[1].
	switch os.Args[1] {
	case "version", "-version", "--version", "-v":
		fmt.Printf("Cetus %s\n", version)
		os.Exit(0)

	case "help", "-help", "--help", "-h":
		// If help was passed then the program shouldn't exit
		// with non-zero error code.
		printUsage()
		os.Exit(0)

	case "set", "fetch":
		// If command & service was not passed then print
		// usage and exit.
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}

	default:
		fmt.Printf("Invalid command: %q\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())

	// If the program has reached this far then that means a valid
	// command was passed & now we should check if a valid service
	// was passed and parse the flags.
	cetus := flag.NewFlagSet("cetus", flag.ExitOnError)

	// We first declare common flags then service specific flags.
	cetus.BoolVar(&dump, "dump", false, "Dump the response")
	cetus.BoolVar(&notify, "notify", false, "Send a desktop notification with info")
	cetus.BoolVar(&print, "print", false, "Print information")
	cetus.BoolVar(&random, "random", false, "Choose a random image")

	switch os.Args[2] {
	case "apod", "nasa":
		defDate := time.Now().UTC().Format("2006-01-02")
		cetus.StringVar(&apodDate, "date", defDate, "Date of NASA APOD to retrieve")
		cetus.Parse(os.Args[3:])

		execAPOD()
	case "bpod", "bing":
		cetus.Parse(os.Args[3:])
		execBPOD()
	default:
		fmt.Printf("Invalid service: %q\n", os.Args[2])
		printUsage()
		os.Exit(1)
	}
}

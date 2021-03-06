package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

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
		defDate := time.Now().UTC().
			// Subtract 8 hours from UTC to ensure program
			// doesn't fail. This still doesn't mean
			// anything, I've emailed them asking about
			// timezone on server but no response :(
			//
			// The server returns "400 Bad Request" when
			// you request future date, so if I request
			// 2020-04-25 on 2020-04-24 23:59 UTC it will
			// return "400 Bad Request" but if I
			// re-request it on 2020-04-25 00:04 UTC it
			// returns "500 Internal Server Error", I
			// think the API server runs on UTC but the
			// program that is responsible for syncing the
			// images is running on a different timezone,
			// which is why it returns "500 Internal
			// Server Error" instead of "400 Bad Request".
			//
			// Hopefully this should work, it will work if
			// the program responsible for syncing images
			// is in or before UTC-8.
			Add(time.Duration(-8) * time.Hour).
			Format("2006-01-02")

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

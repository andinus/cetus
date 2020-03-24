package main

import "fmt"

func printUsage() {
	fmt.Println("Usage: cetus <command> <service> [<flags>]")
	fmt.Println("\nCommands: ")
	fmt.Println(" set     Set the background")
	fmt.Println(" fetch   Fetch the response only")
	fmt.Println(" help    Print help")
	fmt.Println(" version Print Cetus version")
	fmt.Println("\nServices: ")
	fmt.Println(" apod   NASA Astronomy Picture of the Day")
	fmt.Println(" bpod   Bing Photo of the Day")
}

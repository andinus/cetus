package bpod

import (
	"fmt"
)

// Print will print the json output
func Print(res Res) {
	fmt.Printf("Title: %s\n\n", res.Title)
	fmt.Printf("Copyright: %s\n", res.Copyright)
	fmt.Printf("Copyright Link: %s\n", res.CopyrightLink)
	fmt.Printf("Date: %s\n\n", res.StartDate)
	fmt.Printf("URL: %s\n", res.Url)
}

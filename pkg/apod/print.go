package apod

import (
	"fmt"
)

// Print will print the json output
func Print(res Res) {
	fmt.Printf("Title: %s\n\n", res.Title)
	fmt.Printf("Copyright: %s\n", res.Copyright)
	fmt.Printf("Date: %s\n\n", res.Date)
	fmt.Printf("Media Type: %s\n", res.MediaType)
	if res.MediaType == "image" {
		fmt.Printf("URL: %s\n\n", res.HDURL)
	} else {
		fmt.Printf("URL: %s\n\n", res.URL)
	}
	fmt.Printf("Explanation: %s\n", res.Explanation)
}

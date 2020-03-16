package cetus

import (
	"fmt"
	"log"
)

var version string = "v0.4.11"

// Version prints cetus version
func Version() {
	fmt.Printf("Cetus %s\n", version)
}

// PrintPath prints the path passed
func PrintPath(path string) {
	fmt.Println(path)
}

// ErrChk logs the context & error
func ErrChk(ctx string, err error) {
	if err != nil {
		log.Println(ctx)
		log.Fatal(err)
	}
}

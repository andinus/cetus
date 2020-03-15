package cetus

import (
	"fmt"
	"log"
)

var version string = "v0.4.7"

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
	log.Println(ctx)
	log.Fatal(err)
}

package cetus

import "fmt"

var version string = "v0.4.6"

// Version prints cetus version
func Version() {
	fmt.Printf("Cetus %s\n", version)
}

// PrintPath prints the path passed
func PrintPath(path string) {
	fmt.Println(path)
}

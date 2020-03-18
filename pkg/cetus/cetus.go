package cetus

import (
	"fmt"
)

var version string = "v0.5.0"

// Version prints cetus version
func Version() {
	fmt.Printf("Cetus %s\n", version)
}

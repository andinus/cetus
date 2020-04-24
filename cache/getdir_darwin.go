// +build darwin

package cache

import (
	"fmt"
	"os"
)

// GetDir returns cetus cache directory. Default cache directory on
// macOS is $HOME/Library/Caches.
func GetDir() string {
	cacheDir := fmt.Sprintf("%s/%s/%s",
		os.Getenv("HOME"),
		"Library",
		"Caches")

	// Cetus cache directory is cacheDir/cetus
	cetusCacheDir := fmt.Sprintf("%s/%s", cacheDir,
		"cetus")

	return cetusCacheDir
}

// Dir returns "/dev/null", this is required because unveil func in
// main.go calls it & it's useless on macOS anyways so we return
// "/dev/null".
func Dir() string {
	return "/dev/null"
}

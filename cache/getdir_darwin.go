// +build darwin

package cache

import (
	"fmt"
	"os"
)

// GetDir returns cetus cache directory. Default cache directory on
// macOS is $HOME/Library/Caches.
func GetDir() string {
	cacheDir = fmt.Sprintf("%s/%s/%s",
		os.Getenv("HOME"),
		"Library",
		"Caches")

	// Cetus cache directory is cacheDir/cetus
	cetusCacheDir = fmt.Sprintf("%s/%s", cacheDir,
		"cetus")

	return cetusCacheDir
}

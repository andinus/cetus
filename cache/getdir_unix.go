// +build linux netbsd openbsd freebsd dragonfly

package cache

import (
	"fmt"
	"os"
)

// GetDir returns cetus cache directory. Check if the user has set
// CETUS_CACHE_DIR, if not then check if XDG_CACHE_HOME is set & if
// that is not set then assume it to be the default value which is
// $HOME/.cache according to XDG Base Directory Specification.
func GetDir() string {
	cacheDir := Dir()

	// Cetus cache directory is cacheDir/cetus.
	cetusCacheDir := fmt.Sprintf("%s/%s", cacheDir,
		"cetus")

	return cetusCacheDir
}

// Dir returns the system cache directory, this is useful for unveil
// in OpenBSD.
func Dir() string {
	cacheDir := os.Getenv("CETUS_CACHE_DIR")
	if len(cacheDir) == 0 {
		cacheDir = os.Getenv("XDG_CACHE_HOME")
	}
	if len(cacheDir) == 0 {
		cacheDir = fmt.Sprintf("%s/%s", os.Getenv("HOME"),
			".cache")
	}

	return cacheDir
}

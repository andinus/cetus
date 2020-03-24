// +build darwin

package background

import (
	"os/exec"
	"strconv"
)

// Set calls feh to set the background
func Set(path string) error {
	err := exec.Command("osascript", "-e", `tell application "System Events" to tell every desktop to set picture to `+strconv.Quote(path)).Run()
	return err
}

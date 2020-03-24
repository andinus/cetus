// +build darwin

package background

import (
	"fmt"
	"os/exec"
	"strconv"
)

// SetFromFile takes a string as an input, it must be absolute path to
// the background. Checks are not made to check if the path exists or
// it is actually an image, that must be verified before passing it to
// SetFromFile. SetFromFile will exit returning in error if there is
// any.
func SetFromFile(path string) error {
	err := exec.Command("osascript", "-e",
		`tell application "System Events" to tell every desktop to set picture to `+strconv.Quote(path)).Run()
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"set_darwin.go: failed to set background",
			err.Error())
	}
	return err
}

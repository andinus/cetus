// +build linux netbsd openbsd freebsd dragonfly

package background

import (
	"fmt"
	"os"
	"os/exec"
)

// SetFromFile takes a string as an input, it must be absolute path to
// the background. Checks are not made to check if the path exists or
// it is actually an image, that must be verified before passing it to
// SetFromFile. SetFromFile will exit returning in error if there is
// any.
func SetFromFile(path string) error {
	var err error
	switch os.Getenv("XDG_CURRENT_DESKTOP") {
	case "GNOME", "Unity", "Pantheon":
		// GNOME, Unity & Pantheon support setting background
		// from gsettings & have the same key.

		// gsettings takes path in format of a uri
		path = fmt.Sprintf("%s%s", "file://", path)

		err = exec.Command("gsettings",
			"set org.gnome.desktop.background picture-uri", path).Run()
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"set_unix.go: failed to set background with gsettings",
				err.Error())
		}
		return err

	case "LXDE":
		// Background on LXDE can be set with pcmanfm (default
		// file manager).
		err = exec.Command("pcmanfm", "-w", path).Run()
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"set_unix.go: failed to set background with pcmanfm",
				err.Error())
		}
		return err

	default:
		// If WM/DE doesn't have a case then feh is used to
		// set the background. This is tested to work on WMs
		// similar to i3wm.
		feh, err := exec.LookPath("feh")
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"set_unix.go: feh not found in $PATH",
				err.Error())
			return err
		}

		err = exec.Command(feh, "--bg-fill", path).Run()
		if err != nil {
			err = fmt.Errorf("%s\n%s",
				"set_unix.go: failed to set background with feh",
				err.Error())
		}
		return err
	}
}

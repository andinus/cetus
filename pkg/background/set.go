package background

import (
	"fmt"
	"os"
	"os/exec"
)

// Set calls feh to set the background
func Set(path string) error {
	var err error
	switch os.Getenv("XDG_CURRENT_DESKTOP") {
	case "GNOME", "Unity", "Pantheon":
		path = fmt.Sprintf("%s%s", "file://", path)
		err = exec.Command("gsettings", "set org.gnome.desktop.background picture-uri", path).Run()
		return err
	case "LXDE":
		err = exec.Command("pcmanfm", "-w", path).Run()
		return err
	default:
		feh, err := exec.LookPath("feh")
		if err != nil {
			return err
		}
		err = exec.Command(feh, "--bg-fill", path).Run()
		return err
	}
}

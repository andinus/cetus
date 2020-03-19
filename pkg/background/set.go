package background

import "os/exec"

// Set calls feh to set the background
func Set(path string) error {
	feh, err := exec.LookPath("feh")
	if err != nil {
		return err
	}
	err = exec.Command(feh, "--bg-fill", path).Run()
	return err
}

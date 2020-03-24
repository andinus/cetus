// +build linux netbsd openbsd freebsd dragonfly

package notification

import (
	"fmt"
	"os/exec"
)

// Notify sends a desktop notification to the user using libnotify. It
// handles information in the form of Notif struct. It returns an
// error (if exists).
func (n Notif) Notify() error {
	err := exec.Command("notify-send", n.Title, n.Message).Run()
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"notify_unix.go: failed to sent notification with notify-send",
			err.Error())
	}
	return err
}

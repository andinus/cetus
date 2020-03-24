// +build darwin

package notification

import (
	"fmt"
	"os/exec"
	"strconv"
)

// Notify sends a desktop notification to the user using osascript. It
// handles information in the form of Notif struct. It returns an
// error (if exists).
func (n Notif) Notify() error {
	// This script cuts out parts of notification, this bug was
	// confirmed on macOS Catalina 10.15.3, fix not yet known.
	err := exec.Command("osascript", "-e",
		`display notification `+strconv.Quote(n.Message)+` with title `+strconv.Quote(n.Title)).Run()
	if err != nil {
		err = fmt.Errorf("%s\n%s",
			"notify_darwin.go: failed to sent notification with osascript",
			err.Error())
	}
	return err
}

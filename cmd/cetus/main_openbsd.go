// +build openbsd

package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/unix"
	"tildegit.org/andinus/cetus/cache"
)

func main() {
	unveil()
	app()
}

func unveil() {
	unveilL := make(map[string]string)

	// We unveil the whole cache directory.
	err = unix.Unveil(cache.Dir(), "rwc")
	if err != nil {
		log.Fatal(err)
	}

	unveilL["/dev/null"] = "rw" // required by feh
	unveilL["/etc/resolv.conf"] = "r"

	// ktrace output
	unveilL["/usr/libexec/ld.so"] = "r"
	unveilL["/var/run/ld.so.hints"] = "r"
	unveilL["/usr/lib"] = "r"
	unveilL["/dev/urandom"] = "r"
	unveilL["/etc/hosts"] = "r"
	unveilL["/etc/ssl"] = "r"

	for k, v := range unveilL {
		err = unix.Unveil(k, v)
		if err != nil && err.Error() == "no such file or directory" {
			log.Printf("WARN: Unveil failed on %s", k)
		} else if err != nil {
			log.Fatal(fmt.Sprintf("%s :: %s\n%s", k, v,
				err.Error()))
		}
	}

	err = unveilCmd("feh")
	if err != nil {
		log.Fatal(err)
	}

	// Block further unveil calls
	err = unix.UnveilBlock()
	if err != nil {
		log.Fatal(err)
	}
}

// unveilCmd will unveil commands.
func unveilCmd(cmd string) error {
	pathList := strings.Split(getEnv("PATH", ""), ":")
	for _, path := range pathList {
		err = unix.Unveil(fmt.Sprintf("%s/%s", path, cmd), "rx")

		if err != nil && err.Error() != "no such file or directory" {
			return err
		}
	}
	return nil
}

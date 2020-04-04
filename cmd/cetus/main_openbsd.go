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

	unveilL[cache.GetDir()] = "rwc"
	unveilL["/dev/null"] = "rw" // required by feh

	unveilL["/etc/resolv.conf"] = "r"

	// ktrace output
	unveilL["/usr/libexec/ld.so"] = "r"
	unveilL["/var/run/ld.so.hints"] = "r"
	unveilL["/usr/lib/libpthread.so.26.1"] = "r"
	unveilL["/usr/lib/libc.so.95.1"] = "r"
	unveilL["/dev/urandom"] = "r"
	unveilL["/etc/mdns.allow"] = "r"
	unveilL["/etc/hosts"] = "r"
	unveilL["/usr/local/etc/ssl/cert.pem"] = "r"
	unveilL["/etc/ssl/cert.pem"] = "r"
	unveilL["/etc/ssl/certs"] = "r"
	unveilL["/system/etc/security/cacerts"] = "r"
	unveilL["/usr/local/share/certs"] = "r"
	unveilL["/etc/pki/tls/certs"] = "r"
	unveilL["/etc/openssl/certs"] = "r"
	unveilL["/var/ssl/certs"] = "r"

	for k, v := range unveilL {
		err = unix.Unveil(k, v)
		if err != nil && err.Error() != "no such file or directory" {
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

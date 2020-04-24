package main

import (
	"log"

	"tildegit.org/andinus/cetus/cache"
	"tildegit.org/andinus/lynx"
)

func main() {
	unveil()
	app()
}

func unveil() {
	paths := make(map[string]string)

	paths[cache.Dir()] = "rwc"
	paths["/dev/null"] = "rw" // required by feh
	paths["/etc/resolv.conf"] = "r"

	// ktrace output
	paths["/usr/libexec/ld.so"] = "r"
	paths["/var/run/ld.so.hints"] = "r"
	paths["/usr/lib"] = "r"
	paths["/dev/urandom"] = "r"
	paths["/etc/hosts"] = "r"
	paths["/etc/ssl"] = "r"

	err := lynx.UnveilPaths(paths)
	if err != nil {
		log.Fatal(err)
	}

	commands := []string{"feh", "gsettings", "pcmanfm", "notify-send"}

	err = lynx.UnveilCommands(commands)
	if err != nil {
		log.Fatal(err)
	}

	// Block further unveil calls
	err = lynx.UnveilBlock()
	if err != nil {
		log.Fatal(err)
	}
}

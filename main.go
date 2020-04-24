package main

import (
	"log"
	"os"

	"tildegit.org/andinus/cetus/cache"
	"tildegit.org/andinus/lynx"
)

var (
	version string = "v0.6.7"
	dump    bool
	random  bool
	notify  bool
	print   bool

	err     error
	body    string
	file    string
	reqInfo map[string]string

	apodDate string
)

func main() {
	initCetus()

	// Early Check: If command was not passed then print usage and
	// exit. Later command & service both are checked, this check
	// is for version command. If not checked then running cetus
	// without any args will fail because os.Args[1] will panic
	// the program & produce runtime error.
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}

	parseArgs()
}

func initCetus() {
	unveil()
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

// Copyright (c) 2020, Andinus <andinus@inventati.org>

// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.

// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	var (
		imgPath string
		err     error
	)

	// Parse flags passed to program
	flag.StringVar(&imgPath, "img-path", "", "Image to set as wallpaper")
	flag.Parse()

	if len(imgPath) > 0 {
		err = setWall(imgPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func setWall(imgPath string) error {
	var err error

	feh, err := exec.LookPath("feh")
	if err != nil {
		fmt.Println("Error: feh is not in $PATH")
		return err
	}

	err = exec.Command(feh, "--bg-fill", imgPath).Run()
	return err
}

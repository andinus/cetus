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
	"log"

	"framagit.org/andinus/cetus/pkg/unsplash"
)

func main() {
	var (
		photoID string
		random  bool

		width  int
		height int

		err error
	)

	flag.StringVar(&photoID, "photo-id", "", "Unsplash Photo ID to set as background")
	flag.BoolVar(&random, "random", true, "Set a random photo as background")

	flag.IntVar(&width, "width", 1920, "Width of the image")
	flag.IntVar(&height, "height", 1080, "Height of the image")

	flag.Parse()

	if len(photoID) != 0 {
		err = unsplash.SetFromID(photoID, width, height)
		errChk(err)
		return
	}

	if random {
		err = unsplash.SetRandom(width, height)
		errChk(err)
		return
	}
}

func errChk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

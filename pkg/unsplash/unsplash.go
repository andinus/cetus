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

package unsplash

import (
	"fmt"
	"net/http"
	"time"

	"framagit.org/andinus/cetus/pkg"
)

// SetFromID sets background from Unsplash Photo ID
func SetFromID(photoID string, width int, height int) error {
	var path string
	var err error

	path = getPathFromID(photoID)
	path = appendSizeToPath(path, width, height)
	err = background.Set(path)
	return err
}

// SetRandom sets a random photo as background
func SetRandom(width int, height int) error {
	var path string
	var err error

	path, err = getPathRandom(width, height)
	if err != nil {
		return err
	}
	err = background.Set(path)
	return err
}

func getPathFromID(photoID string) string {
	var path string
	path = fmt.Sprintf("%s/%s", "https://source.unsplash.com", photoID)
	return path
}

func getPathRandom(width int, height int) (string, error) {
	var err error
	var path string
	var reqPath string

	client := http.Client{
		Timeout: time.Second * 64,
	}

	reqPath = "https://source.unsplash.com"
	reqPath = appendSizeToPath(reqPath, width, height)

	req, err := http.NewRequest(http.MethodGet, reqPath, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Unsplash Source redirects to the photo
	path = res.Request.URL.String()
	return path, nil
}

func appendSizeToPath(path string, width int, height int) string {
	var size string

	size = fmt.Sprintf("%dx%d", width, height)
	path = fmt.Sprintf("%s/%s", path, size)
	return path
}

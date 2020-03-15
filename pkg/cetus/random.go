package cetus

import (
	"math/rand"
)

// RandAlNum returns random alpha-numeric string of specific length
func RandAlNum(n int) string {
	rand.Seed(time.Now().UnixNano())
	const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

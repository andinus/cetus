package main

import "os"

// getEnv will check if the the key exists, if it does then it'll
// return the value otherwise it will return fallback string.
func getEnv(key, fallback string) string {
	// We use os.LookupEnv instead of using os.GetEnv and checking
	// if the length equals 0 because environment variable can be
	// set and be of length 0. User could've set key="" which
	// means the variable was set but the length is 0. There is no
	// reason why user would want to do this over here though.
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

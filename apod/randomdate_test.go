package apod

import (
	"regexp"
	"testing"
)

// TestRandDate tests the RandDate func. We're not actually testing
// much, there are many cases and we really can't test for all, we'll
// have to run this test a lot of times for that. Even after that it
// doesn't guarantee anything because our test itself is flawed.
func TestRandDate(t *testing.T) {
	date := RandDate()
	re := regexp.MustCompile("((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])")
	if !re.MatchString(date) {
		t.Errorf("Date format is incorrect, got %s, want YYYY-MM-DD.", date)
	}
}

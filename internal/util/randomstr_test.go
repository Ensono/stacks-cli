package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRandomString tests that a 6 character string is returned
func TestRandomString(t *testing.T) {

	// get a 6 character random string
	random := RandomString(6)

	assert.Equal(t, 6, len(random))
}

// TestRandomStringIsDifferent checks to see that running the random string
// function multiple times will generate a new random string each time
func TestRandomStringIsDifferent(t *testing.T) {

	var notunique bool

	// create a slice to hold the different random strings
	var random_list []string

	for i := 0; i < 10; i++ {
		random_list = append(random_list, RandomString(6))
	}

	// ensure that all the strings are unique
	var encountered []string

	for v := range random_list {

		// determine if encountered already has the current item
		// if it has then set the notunique flag and break out of the loop
		if SliceContains(encountered, random_list[v]) {
			notunique = true
			break
		}

		encountered = append(encountered, random_list[v])
	}

	// assert that notunique is false
	assert.Equal(t, false, notunique)
}

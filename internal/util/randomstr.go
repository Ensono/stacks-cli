package util

import (
	crypto_rand "crypto/rand"
	"fmt"
)

// RandomString returns a string  with the specified number of
// characters. Useful for creating temporary filenames etc
// Uses the crypto random generator so as to always give a new string,
// using time as the seed was not fast enough for the tests.
func RandomString(length int) string {
	b := make([]byte, length)

	_, _ = crypto_rand.Read(b[:])

	return fmt.Sprintf("%x", b)[:length]

}

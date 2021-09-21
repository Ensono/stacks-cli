package util

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomString returns a string  with the specified number of
// characters. Useful for creating temporary filenames etc
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

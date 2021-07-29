package helper

import (
	"testing"
)

func TestStacksFlags(t *testing.T) {
	flags := StacksFlags()
	if flags == nil {
		t.Error("No Flags returned")
	}
}

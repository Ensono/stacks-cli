package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// create the test slice with the required values
var slice = []string{
	"fred",
	"bloggs",
}

func TestSliceContainsValue(t *testing.T) {
	assert.Equal(t, true, SliceContains(slice, "fred"))
}

func TestSliceDoesNotContainValue(t *testing.T) {
	assert.Equal(t, false, SliceContains(slice, "newvalue"))
}

package util

import (
	"testing"
)

var expected string = "Hello\nWorld"

func TestTransformCRLFWindows(t *testing.T) {

	result := TransformCRLF("Hello\r\nWorld")

	if result != expected {
		t.Errorf("Carriage return has not been removed")
	}
}

func TestTransformCRLFLinux(t *testing.T) {
	result := TransformCRLF("Hello\nWorld")

	if result != expected {
		t.Errorf("Input text should not have been transformed")
	}
}

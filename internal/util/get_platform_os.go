package util

import "runtime"

// GetPlatformOS is a helper function to the runtime environment
// variable
func GetPlatformOS() string {
	return runtime.GOOS
}

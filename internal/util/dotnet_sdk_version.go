package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Masterminds/semver"
)

// DotnetSDKVersion returns th required Dotnet SDK version from the specified file
func DotnetSDKVersion(file string) (string, string, error) {

	// declare variables
	var err error
	var jsonStr []byte
	var sdkGlobal map[string]interface{}
	var msg string

	// check that the file exists
	if Exists(file) {

		// read the contents of the file into a map
		jsonStr, err = os.ReadFile(file)

		if err != nil {
			return "", "", err
		}
	} else {
		jsonStr = []byte(file)
	}

	json.Unmarshal(jsonStr, &sdkGlobal)

	version, _ := NestedMapLookup(sdkGlobal, "sdk", "version")
	rollForward, _ := NestedMapLookup(sdkGlobal, "sdk", "rollForward")

	versionStr := ""
	if version != nil {
		versionStr = fmt.Sprintf("%v", version)
	}

	// if the rollForward is not empty, then change the version constraint based on the string
	if rollForward != nil {
		switch rollForward {
		case "latestPatch", "latestFeature":

			// as this is the latest patch, remove the patch number from the version string
			sv, _ := semver.NewVersion(versionStr)
			versionStr = fmt.Sprintf("%d.%d.x", sv.Major(), sv.Minor())

			msg = fmt.Sprintf("Rolling forward to the latest patch/feature version of .NET SDK: %s", versionStr)
		}
	}

	return versionStr, msg, err
}

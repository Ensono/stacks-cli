package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// DotnetSDKVersion returns th required Dotnet SDK version from the specified file
func DotnetSDKVersion(file string) (string, error) {

	// declare variables
	var err error
	var jsonStr []byte
	var sdkGlobal map[string]interface{}

	// check that the file exists
	if Exists(file) {

		// read the contents of the file into a map
		jsonStr, err = ioutil.ReadFile(file)

		if err != nil {
			return "", err
		}
	} else {
		jsonStr = []byte(file)
	}

	json.Unmarshal(jsonStr, &sdkGlobal)

	res, _ := NestedMapLookup(sdkGlobal, "sdk", "version")

	result := ""
	if res != nil {
		result = fmt.Sprintf("%v", res)
	}
	return result, err
}

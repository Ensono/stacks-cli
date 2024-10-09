package models

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Masterminds/semver"
)

type Version struct {
	Empty    bool
	Original string
	Major    int
	Minor    int
	Feature  int
	Patch    int

	// private properties that are not exposed outside of the struct
	// the pattern that is to be used to split out the version string
	pattern string

	// the raw version string
	raw string

	// set properties to be used by the methods
	constraint  interface{} // set the required version of the tool
	rollForward string      // this is the string that states which version of .NET is allowed

}

func (v *Version) Init(version string, pattern string) {
	v.pattern = pattern
	v.raw = version

	// set the original version string, based on the pattern
	re := *regexp.MustCompile(v.pattern)
	matches := re.FindStringSubmatch(version)

	// set the original version string
	if len(matches) > 0 {
		v.Original = matches[re.SubexpIndex("version")]
	}
}

func (v *Version) Semver() (bool, error) {
	var err error
	var status bool

	// if the constraint has not been set, set a default value
	if v.constraint == nil {
		v.constraint = "*" // set the default constraint to any version
	}

	// Ensure that the version string can be turned into a semantic version
	// this is done by removing charaters that should not be there
	version := strings.ReplaceAll(v.Original, "_", "")

	// create the semantic version constraint object
	c, err := semver.NewConstraint(v.constraint.(string))
	if err != nil {
		return status, err
	}

	// create the semantic version object
	sv, err := semver.NewVersion(version)
	if err != nil {
		return status, err
	}

	// check to see if the version matches the constraint
	result := c.Check(sv)

	// if the result is false, return an error
	if !result {
		err = fmt.Errorf("version %s does not match the constraint %s", v.Original, v.constraint)
	}

	return result, err
}

func (v *Version) SetSemverConstraint(comparator string) {

	// Update the constraint based on the operator
	v.constraint = comparator
}

func (v *Version) DotNet() (bool, error) {
	var err error
	var status bool
	var constraint Version

	// split the version string into its segments
	err = v.DotNetSplit()
	if err != nil {
		return status, err
	}

	if v.constraint != nil {
		constraint = v.constraint.(Version)
	}

	// switch based on the rollforward property
	if v.rollForward == "" {

		// determine if the currently installed version matches the constraint exactly
		if v.Original == v.constraint.(Version).Original {
			status = true
		}
	} else {
		switch v.rollForward {
		case "latestMajor":
			if v.Major >= constraint.Major {
				status = true
			}
		case "latestMinor":

			// for latest minor all of the previous segments must match
			if v.Major == constraint.Major &&
				v.Minor >= constraint.Minor {
				status = true
			}

		case "latestFeature":

			// for latest feature all of the previous segments must match
			if v.Major == constraint.Major &&
				v.Minor == constraint.Minor &&
				v.Feature >= constraint.Feature {
				status = true
			}

		case "latestPatch":

			// for latest patch all of the other parts of the version must match
			if v.Major == constraint.Major &&
				v.Minor == constraint.Minor &&
				v.Feature == constraint.Feature &&
				v.Patch >= constraint.Patch {
				status = true
			}

		}
	}

	return status, err

}

func (v *Version) DotNetSplit() error {

	var err error

	// create a regular expression object to work with
	re := *regexp.MustCompile(v.pattern)
	matches := re.FindStringSubmatch(v.raw)

	// iterate around all of the matches in the regular expression
	names := re.SubexpNames()

	// create a map of the names and their indexes
	// this clears out any empty groups
	named_groups := make(map[string]int)
	for idx, name := range names {
		if name != "" {
			named_groups[name] = idx
		}
	}

	// determine if the named_groups contain all of the necessary parts,
	// if not then the version is determined to be empty. This means that
	// if one component is missing, then the version is empty
	for _, group := range []string{"major", "minor", "feature", "patch"} {
		if !util.SliceContains(names, group) {
			v.Empty = true
			return err
		}
	}

	for name, idx := range named_groups {

		// skip if no name has been specified
		if name == "version" {
			continue
		}

		// get the value of the match
		value := matches[idx]

		// ensure that the value matches the pattern of integers, before attempting
		// to convert
		re_internal := *regexp.MustCompile(`^\d+$`)
		if re_internal.MatchString(value) {

			// convert the value of the match to an integer and then add
			// the value to the map, using the name
			segment, err := strconv.Atoi(matches[idx])

			if err != nil {
				err = fmt.Errorf("error converting %s segment to integer: %s", name, err)
				return err
			}

			v.Set(name, segment)
		}
	}

	return err
}

func (v *Version) DotNetGlobal(content string) error {
	var err error
	var jsonStr []byte
	var sdkGlobal map[string]interface{}

	// if the content is empty return an error
	if content == "" {
		err = fmt.Errorf("no file or content specified")
		return err
	}

	// Attempt to read the content as a file, or use the content as is
	if util.Exists(content) {

		// read the contents of the file into a map
		jsonStr, err = os.ReadFile(content)

		if err != nil {
			return err
		}
	} else {
		jsonStr = []byte(content)
	}

	// Set the property on the struct
	err = json.Unmarshal(jsonStr, &sdkGlobal)
	if err != nil {
		return err
	}

	// set the properties based on the json object
	// this requires the version to be set as another version object so comparisons can be made
	constraint, err := util.NestedMapLookup(sdkGlobal, "sdk", "version")
	if err != nil {
		return err
	}
	vc := Version{}
	vc.Init(constraint.(string), v.pattern)

	// split the version so the parts can be compared
	err = vc.DotNetSplit()
	if err != nil {
		return err
	}
	v.constraint = vc

	// set the rollforward property
	rollforward, err := util.NestedMapLookup(sdkGlobal, "sdk", "rollForward")
	if err != nil {
		if !strings.Contains(err.Error(), "key not found") {
			return err
		}
		err = nil
	}

	if rollforward != nil {
		v.rollForward = rollforward.(string)
	}

	return err
}

func (v *Version) Set(name string, value interface{}) {
	switch name {
	case "original":
		v.Original = value.(string)
	case "major":
		v.Major = value.(int)
	case "minor":
		v.Minor = value.(int)
	case "feature":
		v.Feature = value.(int)
	case "patch":
		v.Patch = value.(int)
	}
}

func (v *Version) IsEmpty() bool {
	return v.Empty
}

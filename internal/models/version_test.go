package models

import (
	"reflect"
	"testing"
)

func TestVersionInit(t *testing.T) {

	// create a set of test cases and the expected output
	tables := []struct {
		description      string
		version_number   string
		pattern          string
		expectedOriginal string
		expectedRaw      string
	}{
		{
			"Testing version with no named groups",
			"100.98.99",
			`(\d+)\.(\d+)\.(\d+)"`,
			"",
			"100.98.99",
		},
		{
			"Check that the version group sets the original version",
			"100.98.99",
			`(?P<version>(\d+)\.(\d+)\.(\d+))`,
			"100.98.99",
			"100.98.99",
		},
		{
			"Ensure that more complicated strings are handled and the original version is set",
			"v100.98.99",
			`v(?P<version>(\d+)\.(\d+)\.(\d+))`,
			"100.98.99",
			"v100.98.99",
		},
	}

	// iterate around all of the test cases
	for _, table := range tables {

		// create a new version object
		version := Version{}
		version.Init(table.version_number, table.pattern)

		t.Log(table.description)

		// check the version object properties against those that are expected
		if version.Original != table.expectedOriginal {
			t.Errorf("Original version number should be '%s': '%s'", table.expectedOriginal, version.Original)
		}

		if version.raw != table.expectedRaw {
			t.Errorf("Raw version number should be '%s': '%s'", table.expectedRaw, version.raw)
		}
	}
}

func TestVersionDotNetSplit(t *testing.T) {

	// define the version_number that will be used in the tests
	var version_number string = "100.98.99"

	// create the test cases in a table
	tables := []struct {
		description   string
		pattern       string
		err           error
		empty         bool
		expectedParts []int
	}{
		{
			"Check that an empty version, with no subsequent groups, is handled",
			`(?P<version>(\d+)\.(\d+)\.(\d+))`,
			nil,
			true,
			[]int{0, 0, 0, 0},
		},
		{
			"Check that a version with specific groups is handled",
			`(?P<version>(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<feature>[0-9]{1})(?P<patch>[0-9]*))`,
			nil,
			false,
			[]int{100, 98, 9, 9},
		},
		{
			"Ensure that if not all components have been named, the version is empty",
			`(?P<version>(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+))`,
			nil,
			true,
			[]int{0, 0, 0, 0},
		},
	}

	// iterate around the test cases
	for _, table := range tables {

		t.Log(table.description)

		// create the version object
		version := Version{}
		version.Init(version_number, table.pattern)

		// call the split function
		err := version.DotNetSplit()

		if err != table.err {
			t.Errorf("Error should be %s: %s", table.err.Error(), err.Error())
		}

		if version.IsEmpty() != table.empty {
			t.Errorf("Version segments should: %v", table.empty)
		}

		// create a slice of the segments in the version object and compare with the expected parts
		segments := []int{version.Major, version.Minor, version.Feature, version.Patch}
		if !reflect.DeepEqual(segments, table.expectedParts) {
			t.Errorf("Version segments should be the same, expected %v: %v", table.expectedParts, segments)
		}
	}
}

func TestVersionDotNet(t *testing.T) {

	// define the version_number that will be used in the tests
	var pattern string = `(?P<version>(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<feature>[0-9]{1})(?P<patch>[0-9]*))`

	// create a set of test cases and the expected output
	tables := []struct {
		description string
		global      string
		version     string
		expected    bool
	}{
		{
			"Check that version matches global constraint",
			`{"sdk": {"version": "6.0.200"}}`,
			"6.0.200",
			true,
		},
		{
			"Check that version does not match global constraint",
			`{"sdk": {"version": "6.0.200"}}`,
			"6.0.215",
			false,
		},
		{
			"Ensure that latestMajor works",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestMajor"}}`,
			"8.0.200",
			true,
		},
		{
			"Ensure that latestMinor works",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestMinor"}}`,
			"6.1.200",
			true,
		},
		{
			"Ensure that latestFeature works",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestFeature"}}`,
			"6.0.300",
			true,
		},
		{
			"Ensure that latestPatch works",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestPatch"}}`,
			"6.0.250",
			true,
		},
		{
			"Check that a patch constraint does not match",
			`{"sdk": {"version": "6.0.250", "rollForward": "latestPatch"}}`,
			"6.0.225",
			false,
		},
	}

	for _, table := range tables {

		t.Log(table.description)

		// initalise new version object
		version := Version{}
		version.Init(table.version, pattern)

		// set the global sdk details
		err := version.DotNetGlobal(table.global)

		if err != nil {
			t.Error(err.Error())
		}

		actual, _ := version.DotNet()

		if actual != table.expected {
			t.Errorf("Expected %v, got %v", table.expected, actual)
		}
	}
}

func TestVersionDotNetGlobal(t *testing.T) {

	// define the version_number that will be used in the tests
	var version_number string = "100.98.99"
	var pattern string = `(?P<version>(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<feature>[0-9]{1})(?P<patch>[0-9]*))`

	// create a set of test cases and the expected output
	tables := []struct {
		description string
		global      string
		err         string
		constraint  string
		rollforward string
	}{
		{
			"Ensure that a blank global string returns an error",
			``,
			"no file or content specified",
			"",
			"",
		},
		{
			"Check that the constraint is correctly set",
			`{"sdk": {"version": "6.0.200"}}`,
			"",
			"6.0.200",
			"",
		},
		{
			"Check that the rollForward is correctly set",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestPatch"}}`,
			"",
			"6.0.200",
			"latestPatch",
		},
	}

	// iterate around the test cases
	for _, table := range tables {

		t.Log(table.description)

		// create a new version object
		version := Version{}
		version.Init(version_number, pattern)

		// call the global function
		err := version.DotNetGlobal(table.global)

		if err == nil {

			// check that the constraint has been set correctly
			if version.constraint.(Version).raw != table.constraint {
				t.Errorf("Constraint should be %s: %s", table.constraint, version.constraint.(Version).raw)
			}

			// ensure that the rollforward is not set
			if version.rollForward != table.rollforward {
				t.Errorf("Rollforward should be %s: %s", table.rollforward, version.rollForward)

			}
		} else {
			if err.Error() != table.err {
				t.Errorf("Error should be %s: %s", table.err, err.Error())
			}
		}
	}
}

func TestVersionSetSemverConstraint(t *testing.T) {

	var pattern string = `(?P<version>(?P<major>[0-9]*)\.(?P<minor>[0-9]*).(?P<patch>[0-9]*))`
	var version_number string = "100.98.99"

	// create a set of test cases
	tables := []struct {
		description string
		constraint  string
		expected    string
	}{
		{
			"Ensure that the comparator is set correctly, using '~'",
			"~100.98.99",
			"~100.98.99",
		},
		{
			"Ensure that the comparator is set correctly, using '^'",
			"^100.98.99",
			"^100.98.99",
		},
	}

	for _, table := range tables {

		version := Version{}
		version.Init(version_number, pattern)

		// Set the comparator
		version.SetSemverConstraint(table.constraint)

		if version.constraint.(string) != table.expected {
			t.Errorf("Constraint should be %s: %s", table.expected, version.constraint.(string))
		}
	}
}

func TestVersionSemver(t *testing.T) {

	var pattern string = `(?P<version>(?P<major>[0-9]*)\.(?P<minor>[0-9]*).(?P<patch>[0-9]*))`
	var version_number string = "100.98.99"

	// create a set of test cases
	tables := []struct {
		description string
		constraint  string
		expected    bool
	}{
		{
			"Ensure that the version matches the constraint",
			"~100.98.99",
			true,
		},
		{
			"Ensure that the version does not match the constraint",
			"~100.98.100",
			false,
		},
	}

	for _, table := range tables {

		t.Log(table.description)

		version := Version{}
		version.Init(version_number, pattern)
		version.SetSemverConstraint(table.constraint)

		actual, _ := version.Semver()

		if actual != table.expected {
			t.Errorf("Expected %v, got %v", table.expected, actual)
		}
	}
}

func TestVersionSet(t *testing.T) {

	// create a set of test cases that exercise the Set method
	tables := []struct {
		description string
		name        string
		value       interface{}
	}{
		{
			"Check that the major version is set correctly",
			"major",
			100,
		},
		{
			"Check that the minor version is set correctly",
			"minor",
			98,
		},
		{
			"Check that the feature version is set correctly",
			"feature",
			9,
		},
		{
			"Check that the patch version is set correctly",
			"patch",
			8,
		},
	}

	for _, table := range tables {

		// create the version object to work with
		version := Version{}

		// use the set method to set the value
		version.Set(table.name, table.value)

		switch table.name {
		case "major":
			if version.Major != table.value {
				t.Errorf("Major version should be %d: %d", table.value, version.Major)
			}
		case "minor":
			if version.Minor != table.value {
				t.Errorf("Minor version should be %d: %d", table.value, version.Minor)
			}
		case "feature":
			if version.Feature != table.value {
				t.Errorf("Feature version should be %d: %d", table.value, version.Feature)
			}
		case "patch":
			if version.Patch != table.value {
				t.Errorf("Patch version should be %d: %d", table.value, version.Patch)
			}
		}
	}
}

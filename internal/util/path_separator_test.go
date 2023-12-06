package util

import "testing"

func TestNormalisePath(t *testing.T) {

	// create a list of tests to run
	tables := []struct {
		path      string
		separator string
		test      string
		msg       string
	}{
		{
			"/home/russells",
			"/",
			"/home/russells",
			"Path should not be modified. [%s, %s] - %s",
		},
		{
			"/home\\russells",
			"/",
			"/home/russells",
			"Path should be modified to use separator '%s': %s - %s",
		},
		{
			"C:\\users\\rseymour",
			"\\",
			"C:\\users\\rseymour",
			"Path should not be modified. [%s, %s] - %s",
		},
		{
			"C:\\users/rseymour",
			"\\",
			"C:\\users\\rseymour",
			"Path should be modified to use separator '%s': %s - %s",
		},
	}

	// iterate around the table of tests
	for _, table := range tables {

		// run the test
		result := NormalisePath(table.path, table.separator)

		// check the result
		if result != table.test {
			t.Errorf(table.msg, table.separator, table.path, result)
		}
	}
}

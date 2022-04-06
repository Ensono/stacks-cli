package util

import (
	"testing"
)

func TestCheckConnectivity(t *testing.T) {

	// create test table to iterate over
	tables := []struct {
		host string
		test bool
		msg  string
	}{
		{
			"github.com",
			false,
			"Unable to connect to 'github.com'",
		},
		{
			"madeuphost.com.uk",
			true,
			"Unable to connect to 'madeuphost.com.uk'",
		},
	}

	// iterate around the test tables and perform the tests
	for _, table := range tables {

		err := CheckConnectivity(table.host)

		result := err != nil

		if result != table.test {
			t.Error(table.msg)
		}

	}
}

package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// TestNestedMapLookup ensures that a string of JOSN can be turned into a map
// and values retrieved from it
func TestNestedMapLookup(t *testing.T) {

	// create a set of tests to run
	tables := []struct {
		jsonStr   string
		test      string
		testError string
		ks        []string
		msg       string
		msgError  string
	}{
		{
			`{"name": "foo"}`,
			"foo",
			"",
			[]string{"name"},
			"Actual value is not correct: %s != %s",
			"No error was expected",
		},
		{
			`{"name": "foo"}`,
			"",
			"needs at least one key",
			[]string{},
			"No value should be returned",
			"Expected error was not received",
		},
		{
			`{"name": "foo"}`,
			"",
			"key not found",
			[]string{"address"},
			"No value should be returned",
			"Expected error was not received",
		},
		{
			`{"profile": {"firstname": "foo", "lastname": "bar"}}`,
			"bar",
			"",
			[]string{"profile", "lastname"},
			"Actual value is not correct: %s != %s",
			"No error was expected",
		},
	}

	// iterate around the tables
	for _, table := range tables {

		// declare the map that the data needs to be put into
		var data map[string]interface{}
		var result string

		// unmarshal the JSON str into data
		json.Unmarshal([]byte(table.jsonStr), &data)

		// get the result from the lookup in the nested search
		res, err := NestedMapLookup(data, table.ks...)

		if res != nil {
			result = fmt.Sprintf("%v", res)
		}

		if result != table.test {
			t.Errorf(table.msg, result, table.test)
		}

		if err != nil && strings.Contains(err.Error(), table.testError) == false {
			t.Error(table.msgError)
		}
	}
}

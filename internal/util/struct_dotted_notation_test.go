package util

import "testing"

func TestGetFieldBYDottedPath(t *testing.T) {

	// create the struct to work with
	testStruct := struct {
		Person struct {
			FirstName string
			LastName  string
		}
	}{
		Person: struct {
			FirstName string
			LastName  string
		}{
			FirstName: "John",
			LastName:  "Smith",
		},
	}

	// create the table of tests to run
	tables := []struct {
		path string
		test string
		msg  string
	}{
		{
			"Person.FirstName",
			"John",
			"First name should be John: %s",
		},
		{
			"Person.LastName",
			"Smith",
			"Last name should be Smith: %s",
		},
		{
			"person.firstname",
			"John",
			"Last name should be John: %s",
		},
	}

	// iterate around the table of tests
	for _, table := range tables {
		value, _ := GetValueByDottedPath(testStruct, table.path)

		if value != table.test {
			t.Errorf(table.msg, value)
		}
	}
}

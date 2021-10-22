package config

import "testing"

func TestGetMapKey(t *testing.T) {

	// create tables of tests
	tables := []struct {
		f    Framework
		test string
		msg  string
	}{
		{
			Framework{
				Type:   "dotnet",
				Option: "webapi",
			},
			"dotnet_webapi",
			"Output should be `dotnet_webapi`, not `%s`",
		},
		{
			Framework{
				Type:   "java",
				Option: "cqrs",
			},
			"java_cqrs",
			"Output should be `java_cqrs`, not `%s`",
		},
	}

	for _, table := range tables {

		// get the key for this framework type
		res := table.f.GetMapKey()

		if res != table.test {
			t.Errorf(table.msg, res)
		}
	}
}

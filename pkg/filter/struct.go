package filter

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Ensono/stacks-cli/internal/util"
	"gopkg.in/yaml.v2"
)

type yamlSubset struct {
	data map[string]interface{}
}

// NewYamlSubset returns a new yamlSubset object
func New() *yamlSubset {
	return &yamlSubset{
		data: make(map[string]interface{}),
	}
}

func (y *yamlSubset) Filter(data interface{}, filters []string) {

	// create a variable to hold the subset
	subset := make(map[string]interface{})

	for _, filter := range filters {

		// check to see if the specified filter has a value, if not
		// move to the next path
		val, _ := util.GetValueByDottedPath(data, filter)
		if val == "" {
			continue
		}

		subset = y.filterStruct(data, filter, subset)
	}

	// set the object data
	y.data = subset
}

func (y *yamlSubset) Buffer() []byte {

	// create the new buffer
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	defer enc.Close()

	// Encode the data
	err := enc.Encode(y.data)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func (y *yamlSubset) String() string {
	return string(y.Buffer())
}

func (y *yamlSubset) WriteFile(path string, perm uint32) error {

	// ensure that the directory for the file already exists
	basepath := filepath.Dir(path)
	err := util.CreateIfNotExists(basepath, fs.FileMode(perm))
	if err != nil {
		return err
	}

	buf := y.Buffer()
	err = os.WriteFile(path, buf, 0666)

	return err
}

// FilterStruct filters a struct based on a list of filters
func (y *yamlSubset) filterStruct(data interface{}, filter string, subset map[string]interface{}) map[string]interface{} {

	// get the reflect value of the data
	val := reflect.ValueOf(data)

	// split the path into filter_items
	filter_items := strings.Split(filter, ".")

	for _, filter_item := range filter_items {

		var tag string
		var key_name string

		// iterate around the fields in this struct and match the lowercase
		// name, this is so that the tag can be analysed
		t := reflect.TypeOf(data)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if strings.ToLower(field.Name) == filter_item {

				// get the tag for the current field
				tag = field.Tag.Get("yaml")

				break
			}
		}

		// determine the name of the key
		// this is based on either the lowecase name of the property, or if the
		// tag has been set if a value has been set or the first component is '-'
		if tag == "" {
			key_name = filter_item
		} else {

			// split the tag into parts
			// the first part states the name of the key
			//	- if not set then use the component
			//  - if it is '-' skip this field
			tag_parts := strings.Split(tag, ",")
			if tag_parts[0] == "" {
				key_name = filter_item
			} else if tag_parts[0] == "-" {
				break
			} else {
				key_name = tag_parts[0]
			}
		}

		// determine the kind for the current component
		property_value := y.caseInsensitiveFieldByName(val, filter_item)

		// determine the kind of the current component
		kind := property_value.Kind()
		switch kind {
		case reflect.Struct:

			var key_value map[string]interface{}
			if reflect.TypeOf(subset[key_name]) == nil {
				key_value = make(map[string]interface{})
			} else {
				key_value = subset[key_name].(map[string]interface{})
			}

			subset[key_name] = y.filterStruct(property_value.Interface(), strings.Join(filter_items[1:], "."), key_value)
		case reflect.String:
			subset[key_name] = property_value.String()
		}

	}

	return subset
}

func (y *yamlSubset) caseInsensitiveFieldByName(v reflect.Value, name string) reflect.Value {
	name = strings.ToLower(name)
	return v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == name })
}

package main

import (
	"testing"
)

func TestPreflightAsset(t *testing.T) {
	//byte slices
	invalidUTF8 := []byte{0xff, 0xfe, 0xfd}
	validJSON := []byte(`{ "foo": ["bar", "barfoo"] }`)
	invalidJSON := []byte(`foo { "foo": ["bar", "barfoo"] } foo`)
	validYAML := []byte(`"foo": "bar"`)
	invalidYAML := []byte(`"foo":foo:bar: "bar"`)
	multilineYAML := []byte(`"foo":
- "bar"
- "foobar"
- "boofar"
- "roobar"
`)
	multilineYAMLConverted := []byte(`{"foo":["bar","foobar","boofar","roobar"]}`)
	validJSONConverted := []byte(`{ "foo": ["bar", "barfoo"] }`)

	var tests = []struct {
		description string
		data        []byte
		isJSON      bool
		success     bool
	}{
		{"invalid_utf8", invalidUTF8, false, false},
		{"valid_json", validJSON, false, true},
		{"invalid_json", invalidJSON, false, false},
		{"valid_yaml", validYAML, false, true},
		{"invalid_yaml", invalidYAML, false, false},
		{"multiline_yaml", multilineYAML, false, true},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := preflightAsset(&test.data, test.isJSON)
			success := (err == nil)
			if success != test.success {
				expectation := "succeed"
				if test.success == false {
					expectation = "fail"
				}
				t.Errorf("test expected to %s", expectation)
			}
		})
	}

	var comparisons = []struct {
		description string
		dataIn      []byte
		dataOut     []byte
		isJSON      bool
		success     bool
	}{
		{"multiline_yaml_conversion", multilineYAML, multilineYAMLConverted, false, true},
		{"json_conversion", validJSON, validJSONConverted, true, true},
	}

	for _, test := range comparisons {
		t.Run(test.description, func(t *testing.T) {
			err := preflightAsset(&test.dataIn, test.isJSON)
			if err != nil && test.success {
				t.Errorf("test expected to run")
			}

			// quick comparison
			if string(test.dataIn) != string(test.dataOut) {
				t.Errorf("Expected %s to match %s", test.dataIn, test.dataOut)
			}
		})
	}
}

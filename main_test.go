package main

import "testing"

func TestRealMain(t *testing.T) {
	var tests = []struct {
		description string
		schema      string
		args        []string
		success     bool
	}{
		{"valid_schema_valid_document", "testdata/schema.yaml", []string{"testdata/valid.yaml"}, true},
		{"valid_schema_invalid_document", "testdata/schema.yaml", []string{"testdata/invalid.yaml"}, false},
		{"no_schema_valid_document", "", []string{"testdata/valid.yaml"}, true},
		{"no_schema_invalid_document", "", []string{"testdata/broken.yaml"}, false},
		{"no_input", "", []string{}, false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ret := realMain(test.schema, test.args)
			success := ret == 0

			if success != test.success {
				expectation := "succeed"
				if test.success == false {
					expectation = "fail"
				}
				t.Errorf("test expected to %s", expectation)
			}
		})
	}
}

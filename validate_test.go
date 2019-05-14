package main

import (
	"testing"
)

func TestValidateBytes(t *testing.T) {
	schemaJSON := []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Sample schema",
  "type": "object",
  "properties": {
    "foo": {
      "type": "string"
    },
    "bar": {
      "type": "string"
    }
  }
}`)
	validYAML := []byte(`foo: "string-a"
bar: "string-b"
`)
	invalidYAML := []byte(`foo: 35
bar: 70`)
	validJSON := []byte(`{"foo":"string-a","bar":"string-b"}`)
	invalidJSON := []byte(`{"foo":35,"bar":70}`)
	invalidSchemaJSON := []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Sample schema",
  "type": "objectaaa",
  "properties": {
    "foo": {
      "type": "string"
    },
    "bar": {
      "type": "string"
    }
  }
}`)
	invalidSchemaYAML := []byte(`---
"$schema": http://json-schema.org/draft-07/schema#
title: Sample schema
type: object
properties:
  foo:
    type: string
  bar:
    type: string
`)
	brokenYAML := []byte(`a: b: c: d: e: {}`)
	brokenJSON := []byte(`{"foo","bar"}`)

	var tests = []struct {
		description string
                input       []byte
		schema      []byte
		success     bool
	}{
		{"valid_yaml", validYAML, schemaJSON, true},
		{"invalid_yaml", invalidYAML, schemaJSON, false},
		{"valid_json", validJSON, schemaJSON, true},
		{"invalid_json", invalidJSON, schemaJSON, false},
		{"noschema_ok", validYAML, nil, true},
		{"noschema_broken", brokenYAML, nil, false},
		{"invalid_schema_json", validYAML, invalidSchemaJSON, false},
		{"invalid_schema_yaml", validYAML, invalidSchemaYAML, false},
		{"broken_yaml", brokenYAML, schemaJSON, false},
		{"broken_json", brokenJSON, schemaJSON, false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := validateBytes(test.input, test.schema)
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
}

func TestValidateFile(t *testing.T) {
        err := validateFile("testdata/valid.yaml", "testdata/schema.json")
        if err != nil {
                t.Errorf("must accept valid instance document")
        }
        err = validateFile("testdata/invalid.yaml", "testdata/schema.json")
        if err == nil {
                t.Errorf("must reject invalid instance document")
        }
}

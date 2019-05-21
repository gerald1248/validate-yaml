package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	au "github.com/logrusorgru/aurora"
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
	validYAML := "testdata/valid.yaml"
	invalidYAML := "testdata/invalid.yaml"
	validJSON := "testdata/valid.json"
	invalidJSON := "testdata/invalid.json"
	schemaJSON := "testdata/schema.json"
	schemaYAML := "testdata/schema.yaml"
	invalidSchemaJSON := "testdata/invalid-schema.json"
	brokenJSON := "testdata/broken.yaml"

	err := validateFile("testdata/valid.yaml", "testdata/schema.json")
	if err != nil {
		t.Errorf("must accept valid instance document")
	}
	err = validateFile("testdata/invalid.yaml", "testdata/schema.json")
	if err == nil {
		t.Errorf("must reject invalid instance document")
	}
	var tests = []struct {
		description string
		input       string
		schema      string
		success     bool
	}{
		{"valid_yaml", validYAML, schemaJSON, true},
		{"invalid_yaml", invalidYAML, schemaJSON, false},
		{"valid_json", validJSON, schemaJSON, true},
		{"invalid_json", invalidJSON, schemaJSON, false},
		{"noschema_ok", validYAML, "", true},
		{"noschema_broken", brokenJSON, "", false},
		{"schema_yaml", validYAML, schemaYAML, true},
		{"invalid_schema_json", validYAML, invalidSchemaJSON, false},
		{"broken_json", brokenJSON, schemaJSON, false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := validateFile(test.input, test.schema)
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

func TestValidateSTDIN(t *testing.T) {
	var tests = []struct {
		description string
		filename    string
		schema      string
		success     bool
	}{
		{"valid_stdin", "testdata/valid.yaml", "testdata/schema.yaml", true},
		{"invalid_stdin", "testdata/invalid.yaml", "testdata/schema.yaml", false},
		{"broken_stdin", "testdata/broken.yaml", "", false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			file, err := os.Open(test.filename)
			if err != nil {
				log.Println(fmt.Sprintf("can't create mock STDIN stream: %s", au.Bold(err.Error())))
			}
			err = validateSTDIN(file, test.schema)
			success := (err == nil)
			if success != test.success {
				expectation := "succeed"
				if test.success == false {
					expectation = "fail"
				}
				t.Errorf("test expected to %s", expectation)
			}
			file.Close()
		})
	}
}

func TestLoadSchema(t *testing.T) {
	var tests = []struct {
		description string
		schema      string
		success     bool
	}{
		{"valid_schema", "testdata/schema.yaml", true},
		{"broken_schema", "testdata/broken.yaml", false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			_, err := loadSchema(test.schema)
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

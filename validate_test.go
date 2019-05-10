package main

import (
	"testing"
)

func TestValidateBytes(t *testing.T) {
	schema := []byte(`{
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
	valid := []byte(`foo: "string-a"
bar: "string-b"
`)
	invalid := []byte(`foo: 35
bar: 70`)

	err := validateBytes(valid, schema)
	if err != nil {
		t.Errorf("must accept valid instance document")
	}

	err = validateBytes(invalid, schema)
	if err == nil {
		t.Errorf("must reject invalid instance document")
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

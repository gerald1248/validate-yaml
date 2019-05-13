package main

import (
	"testing"
)

func TestPreflightAsset(t *testing.T) {
	//byte slices
	invalidUtf8 := []byte{0xff, 0xfe, 0xfd}
	validJSON := []byte("{ \"foo\": [\"bar\", \"barfoo\"] }")
	invalidJSON := []byte("foo { \"foo\": [\"bar\", \"barfoo\"] } foo")
	validYaml := []byte("\"foo\": \"bar\"")
	invalidYaml := []byte("\"foo\":foo:bar: \"bar\"")
	multilineYaml := []byte(`"foo":
- "bar"
- "foobar"
- "boofar"
- "roobar"
`)
	multilineYamlConverted := []byte("{\"foo\":[\"bar\",\"foobar\",\"boofar\",\"roobar\"]}")

	//expect error
	err := preflightAsset(&invalidUtf8, false)
	if err == nil {
		t.Error("Must reject invalid UTF8")
	}

	err = preflightAsset(&invalidJSON, false)
	if err == nil {
		t.Error("Must reject invalid JSON")
	}

	err = preflightAsset(&invalidYaml, false)
	if err == nil {
		t.Error("Must reject invalid YAML")
	}

	//expect success
	err = preflightAsset(&validYaml, false)
	if err != nil {
		t.Errorf("Must accept valid YAML: %v", err)
	}

	err = preflightAsset(&validJSON, false)
	if err != nil {
		t.Errorf("Must accept valid JSON: %v", err)
	}

	//in-place conversion must match predefined result
	err = preflightAsset(&multilineYaml, false)
	if err != nil {
		t.Errorf("Must accept valid multiline YAML: %v", err)
	}
	if string(multilineYaml) != string(multilineYamlConverted) {
		t.Errorf("Expected %s to match %s", multilineYaml, multilineYamlConverted)
	}
}

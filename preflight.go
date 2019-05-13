package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"unicode/utf8"
)

// ensure YAML as well as JSON can be read
// applies only to file-based processing; the server only accepts JSON
func preflightAsset(a *[]byte, isJSON bool) error {
	if len(*a) == 0 {
		return errors.New("input must not be empty")
	}

	if utf8.Valid(*a) == false {
		return errors.New("input must be valid UTF-8")
	}

	// attempt to parse JSON first
	var any interface{}
	err := json.Unmarshal(*a, &any)

	// input is valid JSON
	if err == nil {
		return nil
	}

	// exit condition: flagged as JSON but error found
	if isJSON {
		return errors.New(fmt.Sprintf("invalid JSON: %s", err.Error()))
	}

	// not JSON
	json, err := yaml.YAMLToJSON(*a)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid YAML: %s", err.Error()))
	}

	// successful conversion
	*a = json

	return nil
}

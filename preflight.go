package main

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"unicode/utf8"

	au "github.com/logrusorgru/aurora"
)

// ensure input is not empty, valid UTF-8, well-formed YAML
func preflightAsset(a *[]byte) error {
	if len(*a) == 0 {
		return errors.New("input must not be empty")
	}

	if utf8.Valid(*a) == false {
		return errors.New("input must be valid UTF-8")
	}

	json, err := yaml.YAMLToJSON(*a)
	if err != nil {
		return errors.New(fmt.Sprintf("input must be well-formed YAML: %s", au.Bold(err.Error())))
	}

	// successful conversion
	*a = json

	return nil
}

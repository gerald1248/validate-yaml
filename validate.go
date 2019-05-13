package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/xeipuuv/gojsonschema"
	"strings"

	au "github.com/logrusorgru/aurora"
)

func validateBytes(bytes []byte, schemabytes []byte) error {

	err := preflightAsset(&bytes, false)
	if err != nil {
		return errors.New(fmt.Sprintf("can't parse input: %s", au.Bold(err.Error())))
	}

	var obj interface{}
	if err = json.Unmarshal(bytes, &obj); err != nil {
		return errors.New(fmt.Sprintf("can't unmarshal data: %s", au.Bold(err.Error())))
	}

	if len(schemabytes) > 0 {
	        schemaLoader := gojsonschema.NewStringLoader(string(schemabytes))
	        documentLoader := gojsonschema.NewStringLoader(string(bytes))

	        result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	        if err != nil {
	                return errors.New(fmt.Sprintf("can't validate JSON: %s\n", au.Bold(err.Error())))
	        }

	        if !result.Valid() {
			var report string
	                for i, desc := range result.Errors() {
				if i > 0 {
					report += "; "
				}
				report += fmt.Sprintf("%s", au.Bold(desc))
	                }
	                return errors.New(fmt.Sprintf("invalid JSON: %s\n", au.Bold(report)))
	        }
	}

	return nil
}

func validateFile(path string, jsonschema string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(fmt.Sprintf("can't read %s: %v", path, au.Bold(err)))
	}

	// schema processing
	var schemabytes []byte
	if len(jsonschema) > 0 {
		log(fmt.Sprintf("Loading schema %s...", au.Bold(jsonschema)))
		schemabytes, err = ioutil.ReadFile(jsonschema)

		schemaIsJSON := strings.HasSuffix(jsonschema, ".json")
		err = preflightAsset(&schemabytes, schemaIsJSON)
		if err != nil {
			return errors.New(fmt.Sprintf("can't parse schema: %s", au.Bold(err.Error())))
		}
	}

	log(fmt.Sprintf("Validating %s...", au.Bold(path)))
	return validateBytes(bytes, schemabytes)
}

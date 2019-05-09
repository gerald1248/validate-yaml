package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/xeipuuv/gojsonschema"

	au "github.com/logrusorgru/aurora"
)

func validateBytes(bytes []byte, schemabytes []byte) error {

	err := preflightAsset(&bytes)
	if err != nil {
		return errors.New(fmt.Sprintf("input failed preflight check: %v", err))
	}

	var obj interface{}
	if err = json.Unmarshal(bytes, &obj); err != nil {
		return errors.New(fmt.Sprintf("can't unmarshal data: %v", err))
	}

	if len(schemabytes) > 0 {
	        schemaLoader := gojsonschema.NewStringLoader(string(schemabytes))
	        documentLoader := gojsonschema.NewStringLoader(string(bytes))

	        result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	        if err != nil {
	                return errors.New(fmt.Sprintf("can't validate JSON: %s\n", err.Error()))
	        }

	        if !result.Valid() {
	                for _, desc := range result.Errors() {
	                        log(fmt.Sprintf("=> %s\n", au.Bold(desc)))
	                }
	                return errors.New(fmt.Sprintf("invalid JSON: %s\n", au.Bold(result.Errors()[0])))
	        }
	}

	return nil
}

func validateFile(path string, jsonschema string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(fmt.Sprintf("can't read %s: %v", path, au.Bold(err)))
	}

	var schemabytes []byte
	if len(jsonschema) > 0 {
		schemabytes, err = ioutil.ReadFile(jsonschema)
	}

	return validateBytes(bytes, schemabytes)
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	au "github.com/logrusorgru/aurora"
)

func validateBytes(bytes []byte, schemabytes []byte) error {

	err := preflightAsset(&bytes, false)
	if err != nil {
		return fmt.Errorf("can't parse input: %s", au.Bold(err.Error()))
	}

	var obj interface{}
	if err = json.Unmarshal(bytes, &obj); err != nil {
		return fmt.Errorf("can't unmarshal data: %s", au.Bold(err.Error()))
	}

	if len(schemabytes) > 0 {
		schemaLoader := gojsonschema.NewStringLoader(string(schemabytes))
		documentLoader := gojsonschema.NewStringLoader(string(bytes))

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return fmt.Errorf("can't validate JSON: %s", au.Bold(err.Error()))
		}

		if !result.Valid() {
			var report string
			for i, desc := range result.Errors() {
				if i > 0 {
					report += "; "
				}
				report += fmt.Sprintf("%s", au.Bold(desc))
			}
			return fmt.Errorf("invalid JSON: %s", au.Bold(report))
		}
	} else {
		log(fmt.Sprintf("%s: checking syntax only", au.Cyan(au.Bold("WARN"))))
	}

	return nil
}

func validateFile(path string, jsonschema string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", path, au.Bold(err))
	}

	schemabytes, err := loadSchema(jsonschema)
	if err != nil {
		return fmt.Errorf("can't parse schema: %s", au.Bold(err.Error()))
	}

	log(fmt.Sprintf("Validating %s...", au.Bold(path)))
	return validateBytes(bytes, schemabytes)
}

func validateSTDIN(jsonschema string) (bool, error) {
	var stdin []byte
	stdinFileInfo, _ := os.Stdin.Stat()
	if stdinFileInfo.Mode()&os.ModeNamedPipe != 0 {
		stdin, _ = ioutil.ReadAll(os.Stdin)
	}

	// empty slice is fine so handle in caller
	if len(stdin) == 0 {
		return true, nil
	}

	schemabytes, err := loadSchema(jsonschema)
	if err != nil {
		return false, fmt.Errorf("can't parse schema: %s", au.Bold(err.Error()))
	}

	log("Validating stream...")
	return false, validateBytes(stdin, schemabytes)
}

func loadSchema(jsonschema string) ([]byte, error) {
	var schemabytes []byte
	var err error
	if len(jsonschema) > 0 {
		log(fmt.Sprintf("Loading schema %s...", au.Bold(jsonschema)))
		schemabytes, err = ioutil.ReadFile(jsonschema)

		schemaIsJSON := strings.HasSuffix(jsonschema, ".json")
		err = preflightAsset(&schemabytes, schemaIsJSON)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't parse schema: %s", au.Bold(err.Error())))
		}
	}

	return schemabytes, nil
}

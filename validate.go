package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Println(fmt.Sprintf("%s: checking syntax only", au.Cyan(au.Bold("WARN"))))
	}

	return nil
}

func validateFile(path string, schema string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", path, au.Bold(err))
	}

	schemabytes, err := loadSchema(schema)
	if err != nil {
		return fmt.Errorf("can't parse schema: %s", au.Bold(err.Error()))
	}

	log.Println(fmt.Sprintf("Validating %s...", au.Bold(path)))
	return validateBytes(bytes, schemabytes)
}

func validateSTDIN(file *os.File, schema string) error {
	var stdin []byte
	stdin, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read input stream: %s", au.Bold(err.Error()))
	}

	schemabytes, err := loadSchema(schema)
	if err != nil {
		return fmt.Errorf("can't parse schema: %s", au.Bold(err.Error()))
	}

	log.Println("Validating stream...")
	return validateBytes(stdin, schemabytes)
}

func loadSchema(schema string) ([]byte, error) {
	var schemabytes []byte
	var err error
	if len(schema) > 0 {
		log.Println(fmt.Sprintf("Loading schema %s...", au.Bold(schema)))
		schemabytes, err = ioutil.ReadFile(schema)

		schemaIsJSON := strings.HasSuffix(schema, ".json")
		err = preflightAsset(&schemabytes, schemaIsJSON)
		if err != nil {
			return nil, fmt.Errorf("can't parse schema: %s", au.Bold(err.Error()))
		}
	}

	return schemabytes, nil
}

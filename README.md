# validate-yaml

`validate-yaml` performs JSON Schema validation for YAML files. It does so by converting the YAML inputs to JSON first.

## Usage
```
$ validate-yaml -h
USAGE validate-yaml [-schema=PATH] FILE [FILE]
  -schema string
    	path to JSON schema definition
$ validate-yaml -schema=schema.json valid.json
OK
```

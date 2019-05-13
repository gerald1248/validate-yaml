# validate-yaml

`validate-yaml` performs JSON Schema validation for YAML files. Under the hood, validation is always against a JSON schema. The YAML input is converted to JSON and then checked against the schema.

In fact, both schema and input files may be YAML documents.

|               |JSON|YAML
|---------------|----|----
|schema document|✅  |✅
|input document |✅  |✅

## Usage
```
$ validate-yaml -h
USAGE validate-yaml [-schema=PATH] FILE [FILE]
  -schema string
    	path to JSON schema definition
$ validate-yaml --schema=schema.json valid.json
OK
```

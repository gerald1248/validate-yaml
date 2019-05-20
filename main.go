package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	au "github.com/logrusorgru/aurora"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"%s %s [-schema=PATH] FILE [FILE]\n",
			au.Cyan(au.Bold("USAGE")),
			au.Bold(filepath.Base(os.Args[0])))
		flag.PrintDefaults()
		os.Exit(0)
	}

	var schema *string
	schema = flag.String("schema", "", "path to JSON schema definition")
	flag.Parse()

	os.Exit(realMain(*schema, flag.Args()))
}

func realMain(schema string, args []string) int {
	errors := 0

	if len(args) == 0 {
		// attempt STDIN handling
		isEmpty, err := validateSTDIN(schema)
		if isEmpty { // special case 1: no input
			log(fmt.Sprintf("No input"))
			flag.Usage()
			return 1
		} else if err != nil { // special case 2: error
			log(fmt.Sprintf("%s %s", au.Red(au.Bold("ERROR")), err.Error()))
			errors++
		}
	} else {
		// iterate over files; return number of invalid files
		for _, arg := range args {
			err := validateFile(arg, schema)
			if err != nil {
				log(fmt.Sprintf("%s %s", au.Red(au.Bold("ERROR")), err.Error()))
				errors++
			}
		}
	}

	// report success only if all files/streams have passed
	if errors == 0 {
		log(fmt.Sprintf("%s", au.Green(au.Bold("OK"))))
	}

	return errors
}

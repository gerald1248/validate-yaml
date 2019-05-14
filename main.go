package main

import (
	"flag"
	"fmt"
	au "github.com/logrusorgru/aurora"
	"os"
	"path/filepath"
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

	var jsonschema *string
	jsonschema = flag.String("schema", "", "path to JSON schema definition")
	flag.Parse()

	errors := 0

	args := flag.Args()
	if len(args) == 0 {
		// attempt STDIN handling
		isEmpty, err := validateSTDIN(*jsonschema)
		if isEmpty { // special case 1: no input
			log(fmt.Sprintf("No input"))
			flag.Usage()
			os.Exit(1)
		} else if err != nil { // special case 2: error
			log(fmt.Sprintf("%s %s", au.Red(au.Bold("ERROR")), err.Error()))
			errors++
		}
	} else {
		// iterate over files; return number of invalid files
		errors := 0
		for _, arg := range args {
			err := validateFile(arg, *jsonschema)
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

	os.Exit(errors)
}

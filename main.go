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

	var jsonschema *string
	jsonschema = flag.String("schema", "", "path to JSON schema definition")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "%s No input files given\n", au.Cyan(au.Bold("INFO")))
		flag.Usage()
		os.Exit(1)
	}

	errors := 0
	for _, arg := range args {
		err := validateFile(arg, *jsonschema)
		if err != nil {
			log(fmt.Sprintf("%s %s", au.Red(au.Bold("ERROR")), err.Error()))
			errors++
		}
	}
	if errors == 0 {
		log(fmt.Sprintf("%s", au.Green(au.Bold("OK"))))
	}

	os.Exit(errors)
}

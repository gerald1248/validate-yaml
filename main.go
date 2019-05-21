package main

import (
	"flag"
	"fmt"
	"log"
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

	// handle possible STDIN input first
	// distinguish from "no input"
	hasNamedPipe := false
	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeNamedPipe != 0 {
		hasNamedPipe = true
	}

	if !hasNamedPipe && len(args) == 0 {
		log.Println("No input")
		flag.Usage()
		return 1
	}

	if hasNamedPipe {
		err := validateSTDIN(os.Stdin, schema)
		if err != nil {
			log.Println(fmt.Sprintf("%s %s", au.Red(au.Bold("ERROR")), err.Error()))
			errors++
		}
	}

	// iterate over files; return number of invalid files
	for _, arg := range args {
		err := validateFile(arg, schema)
		if err != nil {
			log.Println(fmt.Sprintf("%s %s", au.Red(au.Bold("ERROR")), err.Error()))
			errors++
		}
	}

	// report success only if all files/streams have passed
	if errors == 0 {
		log.Println(fmt.Sprintf("%s", au.Green(au.Bold("OK"))))
	}

	return errors
}

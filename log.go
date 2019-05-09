package main

import (
	"fmt"
	"os"
)

func log(s string) {
	fmt.Fprintf(os.Stderr, "%s\n", s)
}

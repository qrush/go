package main

import (
	"os"
	"regexp"
	"fmt"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: regexp [string]")
		os.Exit(1)
	} else if matched, _ := regexp.MatchString("^[0-9]+$", os.Args[1]); matched {
		fmt.Println("Totally numbers")
	} else {
		fmt.Println("Totally not numbers")
	}
}

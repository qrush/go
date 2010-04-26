package main

import (
	"games"
	"rps"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) == 2 {
		file, _ := os.Open(os.Args[1], os.O_RDWR, 0)
		games.Play(os.Stdin, os.Stdout, file, file, rps.NewReferee())
	} else if len(os.Args) == 1 {
		games.Play(os.Stdin, os.Stdout, os.Stdin, os.Stdout, rps.NewReferee())
	} else {
		fmt.Printf("Invalid arguments. Usage: %s [terminal]", os.Args[0])
	}
}

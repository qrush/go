package main

import (
	"games"
	"rps"
	"os"
)

func main() {
	file, _ := os.Open(os.Args[1], os.O_RDWR, 0)
	games.Play(os.Stdin, os.Stdout, file, file, rps.NewReferee())
}

package main

import (
	"os"
	"games"
	"ttt"
)

func main() { games.Play(os.Stdin, os.Stdout, os.Stdin, os.Stdout, ttt.NewReferee()) }

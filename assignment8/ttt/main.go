package main

import (
	"os"
	"games"
	"ttt"
)

func main() { 
	//games.Play(os.Stdin, os.Stdout, os.Stdin, os.Stdout, ttt.NewReferee()) 
	player1 := games.NewLocalView("A", os.Stdin, os.Stdout)
	player2 := games.NewProxyView("B", os.Args[1], os.Args[2])
	ref := ttt.NewReferee()
	for !ref.Turn(player1, player2) {
	}
}

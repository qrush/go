package main

import (
	"flag"
	"games"
	"os"
	"ttt"
)

func main() {
	var player1, player2 games.View
	if *games.PlayerA {
		player1 = games.NewLocalView("A", os.Stdin, os.Stdout)
		player2 = games.NewProxyView("B", flag.Arg(0), flag.Arg(1))
	} else {
		player1 = games.NewProxyView("A", flag.Arg(0), flag.Arg(1))
		player2 = games.NewLocalView("B", os.Stdin, os.Stdout)
	}
	games.Play(player1, player2, ttt.NewReferee())
}

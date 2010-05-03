package main

import (
	"flag"
	"games"
	"os"
	"rps"
	"time"
)

func main() {
	var player1, player2 games.View
	if *games.PlayerA {
		player1 = games.NewLocalView("A", os.Stdin, os.Stdout)
		player2 = games.NewProxyView("B", flag.Arg(0), flag.Arg(1))
	} else {
		player2 = games.NewProxyView("A", flag.Arg(0), flag.Arg(1))
		player1 = games.NewLocalView("B", os.Stdin, os.Stdout)
	}
	ref := rps.NewReferee()
	for !ref.Turn(player1, player2) {
	}
	time.Sleep(.5 * 1e9)
}

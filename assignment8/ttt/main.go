package main

import (
	"os"
	"games"
	"ttt"
	"flag"
	"time"
)

var playerA *bool = flag.Bool("a", false, "set if this process should be player A")

func main() { 
	//games.Play(os.Stdin, os.Stdout, os.Stdin, os.Stdout, ttt.NewReferee()) 
	flag.Parse()	
	var player1, player2 games.View
	if *playerA {
		player1 = games.NewLocalView("A", os.Stdin, os.Stdout)
		player2 = games.NewProxyView("B", flag.Arg(0), flag.Arg(1))
	} else {
		player1 = games.NewProxyView("A", flag.Arg(0), flag.Arg(1))
		player2 = games.NewLocalView("B", os.Stdin, os.Stdout)
	}
	ref := ttt.NewReferee()
	for !ref.Turn(player1, player2) {
	}
	time.Sleep(.5*1e9)
}

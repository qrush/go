package rps

import "games"
import "os"

var (
	player1 games.View
	player2 games.View
)

const (
	rock = "rock"
	paper = "paper"
	scissors = "scissors"
)


func Play(path string) {
	player1 = games.NewLocalView("A", os.Stdin, os.Stdout)

	file, _ := os.Open(path, os.O_RDWR, 0)
	player2 = games.NewLocalView("B", file, file)

	go player1.Enable()
	go player2.Enable()

	p2m := player2.Get()
	p1m := player1.Get()
	player2.Set(p1m)
	player1.Set(p2m)

	player1.Display()
	player2.Display()

	findWinner(p1m.(string), p2m.(string))
} 

func findWinner(p1m, p2m string) {
	if p1m == p2m {
		player1.Done(games.Draw)
		player2.Done(games.Draw)
	} else {
		if p1m == rock {
			if p2m == paper {
				player1.Done(games.Lose)
				player2.Done(games.Win)
			} else {
				player1.Done(games.Win)
				player2.Done(games.Lose)
			}
		} else if p1m == paper {
			if p2m == rock {
				player1.Done(games.Win)
				player2.Done(games.Lose)
			} else {
				player1.Done(games.Lose)
				player2.Done(games.Win)
			}
		} else if p1m == scissors {
			if p2m == rock {
				player1.Done(games.Lose)
				player2.Done(games.Win)
			} else {
				player1.Done(games.Win)
				player2.Done(games.Lose)
			}
		}
	}
}

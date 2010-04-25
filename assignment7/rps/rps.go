package rps

import "games"
import "os"

var (
	player1 games.View
	player2 games.View
)

const (
	rock     = "rock"
	paper    = "paper"
	scissors = "scissors"
)

func isLegal(m interface{}) bool {
	foo := m.(string)
	if foo == rock || foo == paper || foo == scissors {
		return true
	}
	return false
}

func Referee(path string) {
	player1 = games.NewLocalView("A", os.Stdin, os.Stdout)

	file, _ := os.Open(path, os.O_RDWR, 0)
	player2 = games.NewLocalView("B", file, file)

	for {
		p1d := make(chan string)
		p2d := make(chan string)
		f1 := func(v *games.View, c chan string) {
			var m string
			isDone := false
			for ; !isDone; {
				go v.Enable()
				m = v.Get().(string)
				isDone = isLegal(m)
			}
			c <- m
		}

		go f1(&player1, p1d)
		go f1(&player2, p2d)
		
		p1m := <-p1d
		p2m := <-p2d
		player2.Set(p1m)
		player1.Set(p2m)

		player1.Display()
		player2.Display()

		findWinner(p1m, p2m)
	}
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

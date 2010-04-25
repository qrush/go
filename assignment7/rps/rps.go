package rps

import "games"
import "os"

var (
	player1         games.View
	player2         games.View
	rockOutcome     = map[string]games.Outcome{paper: games.Lose, scissors: games.Win}
	scissorsOutcome = map[string]games.Outcome{rock: games.Lose, paper: games.Win}
	paperOutcome    = map[string]games.Outcome{scissors: games.Lose, rock: games.Win}
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
			for !isDone {
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

		Done(p1m, p2m)
	}
}

func Done(p1m, p2m string) {
	p1out := games.Draw
	p2out := games.Draw

	if p1m != p2m {
		switch p1m {
		case rock:
			p1out = rockOutcome[p2m]
		case scissors:
			p1out = scissorsOutcome[p2m]
		case paper:
			p1out = paperOutcome[p2m]
		}

		if p1out == games.Win {
			p2out = games.Lose
		} else {
			p2out = games.Win
		}
	}

	player1.Done(p1out)
	player2.Done(p2out)
}

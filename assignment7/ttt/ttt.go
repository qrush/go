package ttt

import "games"
import "os"
import "math"

var (
	player1 games.View
	player2 games.View
	board   = map[string]int{
		"nw": 0,
		"n":  0,
		"ne": 0,
		"w":  0,
		"c":  0,
		"e":  0,
		"sw": 0,
		"s":  0,
		"se": 0,
	}
	outcomes = [][]string{
		[]string{"nw", "n", "ne"},
		[]string{"w", "c", "e"},
		[]string{"sw", "s", "se"},
		[]string{"nw", "w", "sw"},
		[]string{"n", "c", "s"},
		[]string{"ne", "e", "se"},
		[]string{"nw", "c", "se"},
		[]string{"ne", "c", "sw"},
	}
)

const (
	p1mark = 1
	p2mark = -1
)

func isLegal(m interface{}) bool {
	val, ok := board[m.(string)]
	return ok && val == 0
}

func Referee() {
	player1 = games.NewLocalView("A", os.Stdin, os.Stdout)
	player2 = games.NewLocalView("B", os.Stdin, os.Stdout)

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
		p1m := <-p1d
		Move(p1m, p1mark)
		player2.Set(p1m)
		player2.Display()
		if Done() {
			break
		}

		go f1(&player2, p2d)
		p2m := <-p2d
		Move(p2m, p2mark)
		player1.Set(p2m)
		player1.Display()
		if Done() {
			break
		}
	}
}

func Move(move string, mark int) { board[move] = mark }

func Done() bool {
	done := false
	draw := true

	for _, value := range board {
		if value == 0 {
			draw = false
			break
		}
	}

	if draw {
		player1.Done(games.Draw)
		player2.Done(games.Draw)
		return draw
	}

	for _, outcome := range outcomes {
		first := board[outcome[0]]
		if math.Fabs(float64(first+board[outcome[1]]+board[outcome[2]])) == 3 {
			if first == p1mark {
				player1.Done(games.Win)
				player2.Done(games.Lose)
			} else {
				player1.Done(games.Lose)
				player2.Done(games.Win)
			}

			done = true
		}
	}

	return done
}

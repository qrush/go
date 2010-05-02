package ttt

import (
	"games"
	"math"
	"os"
)

// The current game state, a map of positions => player marks
type tttRef struct {
	board map[string]int
}

// possible winning outcomes for a game
var outcomes = [][]string{
	[]string{"nw", "n", "ne"},
	[]string{"w", "c", "e"},
	[]string{"sw", "s", "se"},
	[]string{"nw", "w", "sw"},
	[]string{"n", "c", "s"},
	[]string{"ne", "e", "se"},
	[]string{"nw", "c", "se"},
	[]string{"ne", "c", "sw"},
}

// A player's mark on the tic tac toe board
const (
	p1mark = 1
	p2mark = -p1mark
)

// Create a new ref for a game with a blank board
func NewReferee() games.Referee {
	return &tttRef{map[string]int{
		"nw": 0,
		"n":  0,
		"ne": 0,
		"w":  0,
		"c":  0,
		"e":  0,
		"sw": 0,
		"s":  0,
		"se": 0,
	}}
}

// Check if a move is in the board and if it hasn't been made
func (this *tttRef) IsLegal(m interface{}) bool {
	val, ok := this.board[m.(string)]
	if m.(string) == "q" {
		os.Exit(0)
	}
	return ok && val == 0
}

// Runs a round of the tic tac toe game.
// * Get move from player 1
// * Update game board
// * Notify player 2
// * Check if the game is done
// * If not, get move from player 2
// * Update game board
// * Notify player 1
// * Check if the game is done
func (this *tttRef) Turn(player1, player2 games.View) (done bool) {
	p1d := make(chan string)
	p2d := make(chan string)

	go games.Listen(this, player1, p1d)
	p1m := <-p1d
	this.board[p1m] = p1mark
	player2.Set(p1m)
	player2.Display()

	if done = this.Done(player1, player2); !done {

		go games.Listen(this, player2, p2d)
		p2m := <-p2d
		this.board[p2m] = p2mark
		player1.Set(p2m)
		player1.Display()
		done = this.Done(player1, player2)

	}

	return
}

// Determines if the game is over. Game is over if:
// * All board elements are marked (draw)
// * A winning outcome is on the board and is owned by one player
func (this *tttRef) Done(player1, player2 games.View) (done bool) {
	draw := true

	for _, value := range this.board {
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
		first := this.board[outcome[0]]
		if math.Fabs(float64(first+this.board[outcome[1]]+this.board[outcome[2]])) == 3 {
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

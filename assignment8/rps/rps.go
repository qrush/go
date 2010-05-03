package rps

import "games"
import "os"

// Game state for RPS, holds both players' current moves.
type rpsRef struct {
	p1move, p2move string
}

// Various outcomes for each throw
var (
	rockOutcome     = map[string]games.Outcome{paper: games.Lose, scissors: games.Win}
	scissorsOutcome = map[string]games.Outcome{rock: games.Lose, paper: games.Win}
	paperOutcome    = map[string]games.Outcome{scissors: games.Lose, rock: games.Win}
)

const (
	rock     = "rock"
	paper    = "paper"
	scissors = "scissors"
)

// Create a new referee for this game
func NewReferee() games.Referee { return &rpsRef{} }

// Determines if this move is legal (is it rock, paper or scissors?)
func (this *rpsRef) IsLegal(m interface{}) bool {
	move := m.(string)
	if move == "q" {
		os.Exit(0)
	}
	return move == rock || move == paper || move == scissors
}

// A round of rock paper scissors:
// * Wait for input from both players at the same time
// * Let the other players know the result
func (this *rpsRef) Turn(player1, player2 games.View) bool {
	p1d := make(chan string)
	p2d := make(chan string)

	go games.Listen(this, player1, p1d)
	go games.Listen(this, player2, p2d)

	this.p1move = <-p1d
	player2.Set(this.p1move)
	this.p2move = <-p2d
	player1.Set(this.p2move)

	player1.Display()
	player2.Display()

	return this.Done(player1, player2)
}

// Checks which player won the game, or if a draw happened.
// * If the moves are the same, draw
// * Otherwise depending on player 1's move
// ** look up the outcome
// ** reverse for the other player
func (this *rpsRef) Done(player1, player2 games.View) bool {
	p1out := games.Draw
	p2out := games.Draw

	if this.p1move != this.p2move {
		switch this.p1move {
		case rock:
			p1out = rockOutcome[this.p2move]
		case scissors:
			p1out = scissorsOutcome[this.p2move]
		case paper:
			p1out = paperOutcome[this.p2move]
		}

		if p1out == games.Win {
			p2out = games.Lose
		} else {
			p2out = games.Win
		}
	}

	player1.Done(p1out)
	player2.Done(p2out)
	return false
}

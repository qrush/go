package rps

import "games"

type rpsRef struct {
	p1move, p2move string
}

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

func NewReferee() games.Referee { return &rpsRef{} }

func (this *rpsRef) IsLegal(m interface{}) bool {
	move := m.(string)
	return move == rock || move == paper || move == scissors
}

func (this *rpsRef) Turn(player1, player2 games.View) bool {
	p1d := make(chan string)
	p2d := make(chan string)

	go games.Listen(this, player1, p1d)
	go games.Listen(this, player2, p2d)

	this.p1move = <-p1d
	this.p2move = <-p2d

	player1.Set(this.p2move)
	player2.Set(this.p1move)

	player1.Display()
	player2.Display()

	return this.Done(player1, player2)
}

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

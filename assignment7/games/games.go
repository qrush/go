package games

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Possible game outcomes
const (
	Win = Outcome(1 + iota)
	Lose
	Draw
)

type (
	// outcome of a game; positive.
	Outcome int

	// what the view must do.
	View interface {
		// enable user interface until view's human player selects a move.
		Enable()

		// return the view's human player's move.
		// If the move is deemed illegal, Enable and Get will be called again.
		Get() interface{}

		// provide other view's human player's move.
		Set(move interface{})

		// reveal other view's human player's move in user interface.
		Display()

		// report game's outcome.
		Done(youWin Outcome)

		// view's event loop, could return a read error.
		Loop() os.Error
	}

	// Implements the View interface, holds state of a player
	LocalView struct {
		gotMove           chan bool
		name              string
		myMove, otherMove string
		in                *bufio.Reader
		out               *bufio.Writer
	}

	// what the referee must doy
	Referee interface {
		// Given the state of two players, is the game over?
		Done(View, View) bool

		// Determines if the given move is acceptable input
		IsLegal(interface{}) bool

		// Represents one round of the game, game continues until this is false
		Turn(View, View) bool
	}
)

// Factory to make a view implementation
func NewLocalView(name string, reader io.Reader, writer io.Writer) (view View) {
	view = new(LocalView)
	view.(*LocalView).gotMove = make(chan bool)
	view.(*LocalView).name = name
	view.(*LocalView).in = bufio.NewReader(reader)
	view.(*LocalView).out = bufio.NewWriter(writer)
	return
}

// Get this view ready for input
func (this *LocalView) Enable() {
	this.Loop()
	this.gotMove <- true
}

// Blocks until a move is made, then returns it
func (this *LocalView) Get() interface{} {
	<-this.gotMove
	return this.myMove
}

// Accepts a move made from another player
func (this *LocalView) Set(move interface{}) { this.otherMove = move.(string) }

// Prints out the player's status
func (this *LocalView) Display() {
	this.out.Write([]byte(fmt.Sprintf("%s's opponent's move: %s\n", this.name, this.otherMove)))
	this.out.Flush()
}

// Report the game outcome to this player
func (this *LocalView) Done(youWin Outcome) {
	switch youWin {
	case Win:
		this.out.Write([]byte(fmt.Sprintf("%s wins.\n", this.name)))
	case Lose:
		this.out.Write([]byte(fmt.Sprintf("%s loses.\n", this.name)))
	case Draw:
		this.out.Write([]byte("Draw game.\n"))
	}
	this.out.Flush()
}

// Accepts input from the player
func (this *LocalView) Loop() os.Error {
	this.out.Write([]byte(fmt.Sprintf("%s's move: ", this.name)))
	this.out.Flush()
	tmp, ok := this.in.ReadBytes('\n')
	this.myMove = string(tmp[0 : len(tmp)-1])
	return ok
}

// Sets up a game with the given readers, writers, and referee
func Play(p1r io.Reader, p1w io.Writer, p2r io.Reader, p2w io.Writer, ref Referee) {
	player1 := NewLocalView("A", p1r, p1w)
	player2 := NewLocalView("B", p2r, p2w)

	for !ref.Turn(player1, player2) {
	}
}

// Repeatedly asks the given view for a legal move on the channel given
func Listen(ref Referee, v View, c chan string) {
	var m string
	isDone := false
	for !isDone {
		go v.Enable()
		m = v.Get().(string)
		isDone = ref.IsLegal(m)
	}
	c <- m
}

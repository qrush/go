package games

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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

	LocalView struct {
		gotMove           chan bool
		name              string
		myMove, otherMove string
		in                *bufio.Reader
		out               *bufio.Writer
	}

	Referee interface {
		Done(View, View) bool
		IsLegal(interface{}) bool
		Turn(View, View) bool
	}
)

func NewLocalView(name string, reader io.Reader, writer io.Writer) (view View) {
	view = new(LocalView)
	view.(*LocalView).gotMove = make(chan bool)
	view.(*LocalView).name = name
	view.(*LocalView).in = bufio.NewReader(reader)
	view.(*LocalView).out = bufio.NewWriter(writer)
	return
}

func (this *LocalView) Enable() {
	this.Loop()
	this.gotMove <- true
}

func (this *LocalView) Get() interface{} {
	<-this.gotMove
	return this.myMove
}

func (this *LocalView) Set(move interface{}) { this.otherMove = move.(string) }

func (this *LocalView) Display() {
	this.out.Write([]byte(fmt.Sprintf("%s's opponent's move: %s\n", this.name, this.otherMove)))
	this.out.Flush()
}

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

func (this *LocalView) Loop() os.Error {
	this.out.Write([]byte(fmt.Sprintf("%s's move: ", this.name)))
	this.out.Flush()
	tmp, ok := this.in.ReadBytes('\n')
	this.myMove = string(tmp[0 : len(tmp)-1])
	return ok
}

func Play(p1r io.Reader, p1w io.Writer, p2r io.Reader, p2w io.Writer, ref Referee) {
	player1 := NewLocalView("A", p1r, p1w)
	player2 := NewLocalView("B", p2r, p2w)

	for !ref.Turn(player1, player2) {
	}
}

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

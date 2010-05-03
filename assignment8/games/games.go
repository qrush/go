package games

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"netchan"
	"os"
	"strings"
	"time"
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

	ProxyView struct {
		gotMove		chan bool
		name		string
		myMove, otherMove	string
		imp		*netchan.Importer
		exp		*netchan.Exporter
		in		chan Move
		out		chan Move
	}

	Move struct {
		m	string
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

var PlayerA *bool = flag.Bool("a", false, "set if this process should be player A")

func init() {
	flag.Parse()
}

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

// Factory to make a proxy view implementation
func NewProxyView(name, local, remote string) (view View) {
	ls := strings.Split(local, ":", 0)
	rs := strings.Split(remote, ":", 0)

	view = new(ProxyView)
	view.(*ProxyView).name    = name
	view.(*ProxyView).exp, _  = netchan.NewExporter("tcp", local)
	view.(*ProxyView).gotMove = make(chan bool)
	view.(*ProxyView).in      = make(chan Move)
	view.(*ProxyView).out     = make(chan Move)

	view.(*ProxyView).exp.Export(ls[1], view.(*ProxyView).out, netchan.Send, new(Move))
	for {
		var err os.Error
		if view.(*ProxyView).imp, err = netchan.NewImporter("tcp", remote); err != nil {
			time.Sleep(1*1e9)
		} else {
			break
		}
	}
	view.(*ProxyView).imp.Import(rs[1], view.(*ProxyView).in, netchan.Recv, new(Move))
	return
}

// Get this view ready for input
func (this *ProxyView) Enable() {
	this.Loop()
	this.gotMove <- true
}

// Blocks until a move is made, then returns it
func (this *ProxyView) Get() interface{} {
	<-this.gotMove
	return this.myMove
}

// Accepts a move made from another player
func (this *ProxyView) Set(move interface{}) {
	this.otherMove = move.(string)
	fmt.Println("SET: SENDING OUT MOVE")
	this.out <- Move{move.(string)}
	fmt.Println("SET: SENT MOVE")
}

// Prints out the player's status
func (this *ProxyView) Display() {
}

// Report the game outcome to this player
func (this *ProxyView) Done(youWin Outcome) {
}

// Accepts input from the player
func (this *ProxyView) Loop() os.Error {
	var tmp Move
	fmt.Println("LOOP: WAITING FOR MOVE")
	tmp = <-this.in
	fmt.Println("LOOP: GOT FOR MOVE")
	this.myMove = tmp.m
	return nil
}


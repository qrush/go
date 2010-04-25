package games

import "os"
import "io"
import "bufio"
import "fmt"

// outcome of a game; positive.
type Outcome int

const (
	Win = Outcome(1 + iota)
	Lose
	Draw
)

// what the view must do.
type View interface {
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

type LocalView struct {
	gotMove           chan bool
	name              string
	myMove, otherMove string
	in                *bufio.Reader
	out               *bufio.Writer
}

func NewLocalView(name string, reader io.Reader, writer io.Writer) View {
	var view View
	view = new(LocalView)
	view.(*LocalView).gotMove = make(chan bool)
	view.(*LocalView).name = name
	view.(*LocalView).in = bufio.NewReader(reader)
	view.(*LocalView).out = bufio.NewWriter(writer)
	return view
}

func (this *LocalView) Enable() {
	this.Loop()
	this.gotMove <- true
}

func (this *LocalView) Get() interface{} {
	<-this.gotMove
	foo := this.myMove
	return foo
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
	tmp = tmp[0 : len(tmp)-1]
	this.myMove = string(tmp)
	return ok
}

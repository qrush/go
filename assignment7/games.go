package games

import "os"

// outcome of a game; positive.
type Outcome int

const (
  Win  = Outcome(1+iota)
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

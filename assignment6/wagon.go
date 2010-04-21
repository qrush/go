/*
 * Everything goes in main because this is a program, not a package.
 */
package main

import (
	"container/list"
	"exec"
	"fmt"
	"os"
	"os/signal"
	"strconv"
)

var (
	width, height int // width and height of the "prairie"
	train         *list.List // the wagon train
	nextDisplay   byte // the display character for the next train to be added
)

type (
	Wagon struct {
		x, y    int
		display string // This is what gets printed when we display the train
	}
)

/* Redraw all the wagons and an info message */
func redraw() {
	clearScreen()
	for w := range train.Iter() {
		this := w.(*Wagon)
		drawAt(this.x, this.y, this.display)
	}
	drawAt(0, height + 1, "udlr to move head, UDLR to move tail, \na to add new head, A to add new tail, q to quit\n")
}

func clearScreen() { fmt.Printf("\033[2J\n") }

/* Draw a string at a given x/y coordinate */
func drawAt(x, y int, s string) { fmt.Printf("\033[%d;%dH%s", y, x, s) }

/* Get the front of the train */
func head() *list.Element { return train.Front() }

/* Get the end of the train */
func tail() *list.Element { return train.Back() }

/* Return the Wagon info for a list element */
func wagon(e *list.Element) *Wagon { return e.Value.(*Wagon) }

/* Move the front wagon and have the rest follow */
func moveFront(x, y int) {
	move(x, y, tail, func(e *list.Element) *list.Element { return e.Prev() })
}

/* Move the tail and have the rest follow */
func moveBack(x, y int) {
	move(x, y, head, func(e *list.Element) *list.Element { return e.Next() })
}

/* Generic moving */
func move(x, y int, start func() *list.Element, advance func(*list.Element) *list.Element) {
	if x > 0 && y > 0 && x <= width && y <= height {
		first := wagon(start())

		for e := start(); e != nil; e = advance(e) {
			if wagon(e) != first {
				first.x = wagon(e).x
				first.y = wagon(e).y
			}
			first = wagon(e)
		}

		first.x = x
		first.y = y
	}
}

/* Process the keystroke from the user */
func process(input string) {
	head := wagon(head())
	tail := wagon(tail())

	switch input {
	case "a":
		addToTrain(true)
	case "u":
		moveFront(head.x, head.y-1)
	case "d":
		moveFront(head.x, head.y+1)
	case "l":
		moveFront(head.x-1, head.y)
	case "r":
		moveFront(head.x+1, head.y)
	case "A":
		addToTrain(false)
	case "U":
		moveBack(tail.x, tail.y-1)
	case "D":
		moveBack(tail.x, tail.y+1)
	case "L":
		moveBack(tail.x-1, tail.y)
	case "R":
		moveBack(tail.x+1, tail.y)
	case "q":
		cleanup()
	}
}

/* Add a new wagon to either the front or the back based on the value of the argument */
func addToTrain(front bool) {
	if nextDisplay < '~' {
		if front {
			train.PushFront(New(1, 1))
		} else {
			train.PushBack(New(width, height))
		}
	}
}

/* Run stty with the given argument string */
func stty(mode string) {
	cmd, err := exec.Run("/bin/stty", []string{"stty", mode}, os.Environ(), "", exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	cmd.Close()
}

/* Make a new Wagon */
func New(x, y int) *Wagon {
	w := new(Wagon)
	w.x = x
	w.y = y
	w.display = string(nextDisplay)
	nextDisplay++
	return w
}

/* Fix tty and exit. */
func cleanup() {
	stty("sane")
	os.Exit(0)
}

/* Set up some default values and add the initial wagons */
func init() {
	width = 30
	height = 30
	if len(os.Args) == 3 {
		width,_ = strconv.Atoi(os.Args[1])
		height,_ = strconv.Atoi(os.Args[2])
	} else if len(os.Args) != 1 {
		fmt.Printf("Usage: %s [width] [height]\n", os.Args[0]);
		os.Exit(1)
	}
	nextDisplay = 'a'
	train = list.New()
	train.PushFront(New(1, 1))
	train.PushBack(New(width, height))

	go func() {
		for {
			if (<-signal.Incoming).(signal.UnixSignal) == 2 {
				cleanup()
			}
		}
	}()
}

/* Kick it all off */
func main() {
	stty("cbreak")
	b := make([]byte, 1)

	for {
		redraw()
		os.Stdin.Read(b)
		fmt.Printf("%s", b)
		process(string(b))
	}
}

package main

import (
	"container/list"
	"exec"
	"fmt"
	"os"
)

var (
	width, height int
	train         *list.List
	nextDisplay   byte
)

type (
	Wagon struct {
		x, y    int
		display string
	}
)

func redraw() {
	clearScreen()
	for w := range train.Iter() {
		this := w.(*Wagon)
		drawAt(this.x, this.y, this.display)
	}
}

func clearScreen() { fmt.Printf("\033[2J\n") }

func drawAt(x, y int, s string) { fmt.Printf("\033[%d;%dH%s\n", y, x, s) }

func head() *list.Element { return train.Front() }

func tail() *list.Element { return train.Back() }

func wagon(e *list.Element) *Wagon { return e.Value.(*Wagon) }

func moveFront(x, y int) {
	move(x, y, tail, func(e *list.Element) *list.Element { return e.Prev() })
}

func moveBack(x, y int) {
	move(x, y, head, func(e *list.Element) *list.Element { return e.Next() })
}

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
		stty("sane")
		os.Exit(0)
	}
}

func addToTrain(front bool) {
	if nextDisplay < '~' {
		if front {
			train.PushFront(New(1, 1))
		} else {
			train.PushBack(New(width, height))
		}
	}
}

func stty(mode string) {
	cmd, err := exec.Run("/bin/stty", []string{"stty", mode}, os.Environ(), "", exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	cmd.Close()
}

func New(x, y int) *Wagon {
	w := new(Wagon)
	w.x = x
	w.y = y
	w.display = string(nextDisplay)
	nextDisplay++
	return w
}

func init() {
	width = 30
	height = 30
	nextDisplay = 'a'
	train = list.New()
	train.PushFront(New(1, 1))
	train.PushBack(New(width, height))
}

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

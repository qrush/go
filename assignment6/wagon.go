package main

import (
	"fmt"
	"container/list"
	"exec"
	"os"
)

var (
	width, height int
	train         *list.List
)

type Wagon struct {
	x, y    int
	display string
}

func redraw() {
	clearScreen()
	for w := range train.Iter() {
		this := w.(*Wagon)
		drawAt(this.x, this.y, this.display)
	}
}

func clearScreen() { fmt.Printf("\033[2J\n") }

func drawAt(x, y int, s string) { fmt.Printf("\033[%d;%dH%s\n", y, x, s) }

func moveFront(x, y int) {
	var first, prev *Wagon
	if x > 0 && y > 0 && x <= width && y <= height {
		iter := train.Iter()
		first = (<-iter).(*Wagon)
		prev = first
		for w := range iter {
			this := w.(*Wagon)
			this.x = prev.x
			this.y = prev.y
			prev = this
		}
		first.x = x
		first.y = y
	}
}

func moveBack(x, y int) {
	var last, next *Wagon
	last = train.Back().Value.(*Wagon)
	next = last

	if x > 0 && y > 0 && x <= width && y <= height {
		for e := train.Back(); e != nil; e = e.Prev() {
			if e.Value.(*Wagon) != last {
				e.Value.(*Wagon).x = next.x
				e.Value.(*Wagon).y = next.y
			}
			next = e.Value.(*Wagon)
		}

		last.x = x
		last.y = y
	}
}

func process(input string) {
	head := (train.Front()).Value.(*Wagon)
	tail := (train.Back()).Value.(*Wagon)

	switch input {
	case "u":
		moveFront(head.x, head.y-1)
	case "d":
		moveFront(head.x, head.y+1)
	case "l":
		moveFront(head.x-1, head.y)
	case "r":
		moveFront(head.x+1, head.y)
	case "U":
		moveBack(tail.x, tail.y-1)
	case "D":
		moveBack(tail.x, tail.y+1)
	case "L":
		moveBack(tail.x-1, tail.y)
	case "R":
		moveBack(tail.x+1, tail.y)
	}

}

func NewWagon(x, y int, display string) *Wagon {
	w := new(Wagon)
	w.x = x
	w.y = y
	w.display = display
	return w
}

func init() {
	width = 30
	height = 30
	train = list.New()
	train.PushFront(NewWagon(1, 1, "a"))
	train.PushBack(NewWagon(width, height, "b"))
}

func main() {
	cmd, err := exec.Run("/bin/stty", []string{"stty", "cbreak"}, os.Environ(), "", exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	cmd.Close()

	b := make([]byte, 1)
	for {
		redraw()
		os.Stdin.Read(b)
		fmt.Printf("%s", b)
		process(string(b))
	}
}

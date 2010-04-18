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
	nextDisplay   byte
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

func moveBack(x, y int) {
	var prev *Wagon
	prev = train.Front().Value.(*Wagon)

	move(x, y, func() {
		for e := train.Front(); e != nil; e = e.Next() {
			if e.Value.(*Wagon) != prev {
				prev.x = e.Value.(*Wagon).x
				prev.y = e.Value.(*Wagon).y
			}
			prev = e.Value.(*Wagon)
		}

		prev.x = x
		prev.y = y
	})
}

func moveFront(x, y int) {
	var next *Wagon
	next = train.Back().Value.(*Wagon)

	move(x, y, func() {
		for e := train.Back(); e != nil; e = e.Prev() {
			if e.Value.(*Wagon) != next {
				next.x = e.Value.(*Wagon).x
				next.y = e.Value.(*Wagon).y
			}
			next = e.Value.(*Wagon)
		}

		next.x = x
		next.y = y
	})
}

func move(x, y int, fn func()) {
	if x > 0 && y > 0 && x <= width && y <= height {
		fn()
	}
}

func process(input string) {
	head := (train.Front()).Value.(*Wagon)
	tail := (train.Back()).Value.(*Wagon)

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
		os.Exit(0)
	}
}

func addToTrain(front bool) {
	if nextDisplay < '~' {
		if front {
			train.PushFront(NewWagon(1, 1))
		} else {
			train.PushBack(NewWagon(width, height))
		}
	}
}

func cleanup() {
	cmd, err := exec.Run("/bin/stty", []string{"stty", "sane"}, os.Environ(), "", exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	cmd.Close()
}

func NewWagon(x, y int) *Wagon {
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
	train.PushFront(NewWagon(1, 1))
	train.PushBack(NewWagon(width, height))
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

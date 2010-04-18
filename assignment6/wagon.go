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

func drawAt(x, y int, s string) { fmt.Printf("\033[%d;%dH%s\n", x, y, s) }

func process(input string) {}

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
		process(string(b))
		redraw()
		os.Stdin.Read(b)
		fmt.Printf("%s", b)
	}
}

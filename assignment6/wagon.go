package main

import "fmt"
import "time"

func redraw() {
}

func clearScreen() {
	fmt.Printf("\033[2J\n")
}

func drawAt(x, y int, s string) {
	fmt.Printf("\033[%d;%dH%s\n", x, y, s)
}

func main() { 
	clearScreen()
	drawAt(2,3,"A")
	time.Sleep(30000000000)
}

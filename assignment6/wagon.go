package main

import "fmt"
import "container/list"
import "exec"
import "os"

type Wagon struct {
	x, y    int
	display string
}

func redraw() {
}

func clearScreen() {
	fmt.Printf("\033[2J\n")
}

func drawAt(x, y int, s string) {
	fmt.Printf("\033[%d;%dH%s\n", x, y, s)
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
		os.Stdin.Read(b)
		fmt.Printf("%s",b)
	}

	clearScreen()
	drawAt(2,3,"A")

	fmt.Println("WAGON TRAIN")

	l := list.New()
	w1 := new(Wagon)
	w2 := new(Wagon)
	w3 := new(Wagon)

	w1.x = 3
	w1.y = 3
	w1.display = "a"
	w2.x = 4
	w2.y = 4
	w2.display = "b"
	w3.x = 5
	w3.y = 5
	w3.display = "c"

	l.PushBack(w1)
	l.PushBack(w2)
	l.PushBack(w3)

	for w := range l.Iter() {
		fmt.Println(w.(*Wagon).display)
	}
}

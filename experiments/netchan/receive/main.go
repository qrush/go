package main

import "fmt"
import "netchan"
import "os"

func main() {
	fmt.Println("I WILL FIND MY PARENTS")
	imp,ierr := netchan.NewImporter("tcp", ":9292")
	fmt.Println(ierr)
	c := make(chan os.Error)
	err := imp.Import("octocator", c, netchan.Recv, new(os.Error))
	fmt.Println(err)

	if err == nil {
		cat := <-c
		fmt.Println(cat)
	}
}

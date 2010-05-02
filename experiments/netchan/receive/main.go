package main

import "fmt"
import "netchan"

type Octocat struct {
	parents int
}

func main() {
	fmt.Println("I WILL FIND MY PARENTS")
	imp, _ := netchan.NewImporter("tcp", ":9292")
	c := make(chan Octocat)
	err := imp.Import("octocator", c, netchan.Recv, new(Octocat))
	fmt.Println(err)

	if err == nil {
		cat := <-c
		fmt.Println(cat.parents)
	}
}

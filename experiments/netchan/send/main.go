package main

import "fmt"
import "netchan"

type Octocat struct {
	parents int
}

func main() {
	fmt.Println("LOL")
	exp, _ := netchan.NewExporter("tcp", ":9292")
	c := make(chan Octocat)
	err := exp.Export("octocator", c, netchan.Send, new(Octocat))
	fmt.Println(err)
	c <- Octocat{0}
}

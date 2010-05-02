package main

import "fmt"
import "netchan"

type Octocat struct {
	parents int
}

func main() {
	fmt.Println("LOL")
	exp, err := netchan.NewExporter("tcp", ":8383")
	c := make(chan Octocat)
	err2 := exp.Export("name", c, netchan.Send)
	c <- Octocat{0}
}

package main

import "fmt"
import "netchan"

type Octocat struct {
	parents int
}

func main() {
	fmt.Println("LOL")
	exp, _ := netchan.NewExporter("tcp", ":8383")
	c := make(chan Octocat)
	exp.Export("name", c, netchan.Send, new(Octocat))
	c <- Octocat{0}
}

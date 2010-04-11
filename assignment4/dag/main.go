package main

import (
	"dag"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if err := dag.Main(dag.NewTarget, dag.Print); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

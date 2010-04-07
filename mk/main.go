package main

import (
	"dag"
	"flag"
	"fmt"
	"os"
	"mk"
)

func main() {
	flag.Parse()
	if err := dag.Main(mk.NewTarget, mk.Print); err != nil {
	  fmt.Println(err)
	  os.Exit(1)
	}
}

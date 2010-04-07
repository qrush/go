package main

import (
	"flag"
	"fmt"
)

var code *int = flag.Int("areacode", 716, "some other area code")

func main() {
	fmt.Printf("Testing out flags!\n")
	flag.Parse()
	fmt.Println("areacode has value ", *code)
}

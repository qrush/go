package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	bytes, _ := ioutil.ReadFile("example1.ls")
	strs := strings.Fields(string(bytes))

	for _, str := range strs {
		fmt.Println(str)
	}
}

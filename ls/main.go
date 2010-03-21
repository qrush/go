package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
)

type (
	// Scanner function, must return next token or false.
	// Advance past next token if argument is true.
	Scanner func(bool) (string, bool)

	// Parser function
	Parser func(Scanner) Eval

	// Interpreter function; returns value or os.Error.
	Eval func() interface{}

	Node struct {
		Name      string
		Directory *os.Dir
	}
)

var strs string

func makeArr(f func(), funcs []func()) []func() {
	newFuncs := make([]func(), len(funcs)+1)
	newFuncs[0] = f
	for i := range funcs {
		newFuncs[i+1] = funcs[i]
	}
	return newFuncs
}

func name(c Node) func() { return func() { fmt.Println(c.Name) } }

func nl() func() { return func() { fmt.Println("\n") } }

func Expr(c Node, next Scanner) []func() {
	nt, _ := next(true)
	switch nt {
	case "(":
		return Expr(c, next)
	case "(name)":
		return makeArr(name(c), Expr(c, next))
	case "(nl)":
		return makeArr(nl(), Expr(c, next))
	case ")":
		return []func(){}
	}
	return nil
}

func main() {
	bytes, _ := ioutil.ReadFile("example2.ls")
	strs := strings.Fields(string(bytes))

	arg := 0 // consider os.Args[++arg] next

	node := new(Node)
	node.Name = "lolz"
	node.Directory, _ = os.Stat(".")

	result := Expr(*node, func(use bool) (string, bool) {
		switch {
		case arg >= len(strs):
			return "", false
		case use:
			ret := strs[arg]
			arg++
			return ret, true
		}
		return strs[arg], true
	})

	//for _, str := range strs {
	fmt.Println(strs)
	fmt.Println(result)
	//}

	for fn := range result {
		result[fn]()
	}
}

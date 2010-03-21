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
		Name string
		Dir  *os.Dir
	}
)

func makeArr(f func(), funcs []func()) []func() {
	newFuncs := make([]func(), len(funcs)+1)
	newFuncs[0] = f
	for i := range funcs {
		newFuncs[i+1] = funcs[i]
	}
	return newFuncs
}

func name(c Node) func() { return func() { fmt.Printf("%s", c.Name) } }

func nl() func() { return func() { fmt.Println() } }

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
	dname := "."
	fpt, _ := os.Open(dname, os.O_RDONLY, 0666)
	names, _ := fpt.Readdirnames(-1)
	nodes := make([]Node, len(names))
	for i, cn := range names {
		nodes[i].Name = cn
		nodes[i].Dir, _ = os.Stat(cn)
	}

	for _, node := range nodes {
		arg := 0
		strs := strings.Fields(string(bytes))
		scanner := func(use bool) (string, bool) {
			switch {
			case arg >= len(strs):
				return "", false
			case use:
				ret := strs[arg]
				arg++
				return ret, true
			}
			return strs[arg], true
		}

		result := Expr(node, scanner)
		for _, fn := range result {
			fn()
		}
	}
}

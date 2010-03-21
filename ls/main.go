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

	Node struct {
		Name string
		Dir  *os.Dir
	}
)

func compose(f func(), funcs []func()) []func() {
	newFuncs := make([]func(), len(funcs)+1)
	newFuncs[0] = f
	for i := range funcs {
		newFuncs[i+1] = funcs[i]
	}
	return newFuncs
}

func name(n Node) func() { return func() { fmt.Printf("%s", n.Name) } }

func nl() func() { return func() { fmt.Println() } }

func file(n Node, funcs []func()) func() {
	return func() {
		if n.Dir.IsRegular() {
			for _, fn := range funcs {
				fn()
			}
		}
	}
}

func dir(n Node, funcs []func()) func() {
	return func() {
		if n.Dir.IsDirectory() {
			for _, fn := range funcs {
				fn()
			}
		}
	}
}

func Expr(n Node, next Scanner) []func() {
	nt, _ := next(true)
	switch nt {
	case "(":
		return Expr(n, next)
	case "(file":
		return compose(file(n, Expr(n, next)), Expr(n, next))
	case "(dir":
		return compose(dir(n, Expr(n, next)), Expr(n, next))
	case "(name)":
		return compose(name(n), Expr(n, next))
	case "(nl)":
		return compose(nl(), Expr(n, next))
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

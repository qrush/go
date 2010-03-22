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
	Scanner func(bool) (string, bool, int)

	Node struct {
		Name string
		Dir  *os.Dir
	}
)

var bytes []byte

func compose(f func(), funcs []func()) []func() {
	newFuncs := make([]func(), len(funcs)+1)
	newFuncs[0] = f
	for i := range funcs {
		newFuncs[i+1] = funcs[i]
	}
	return newFuncs
}

func compose2(f1 []func(), f2 []func()) []func() {
	newFuncs := make([]func(), len(f1)+len(f2))
	for i := range f1 {
		newFuncs[i] = f1[i]
	}
	for i := range f2 {
		newFuncs[i + len(f1)] = f2[i]
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
/*
	if n.Dir.IsDirectory() {
		return funcs
	}
	return []func() {}
*/
}

func subNodes(n Node) []Node {
	if (n.Dir.IsDirectory()) {
		fpt, _ := os.Open(n.Name, os.O_RDONLY, 0666)
		names, _ := fpt.Readdirnames(-1)
		nodes := make([]Node, len(names))
		for i, cn := range names {
			nodes[i].Name = n.Name + "/" + cn
			nodes[i].Dir, _ = os.Stat(nodes[i].Name)
			//fmt.Printf("Filename: %s, parent dir: %s\n", nodes[i].Name, n.Name)
		}
		return nodes
	} 
	return []Node {}
}

func sub(n Node, funcs []func()) func() {
	return func() {
		if n.Dir.IsDirectory() {
			for _, fn := range funcs {
				fn()
			}
		}
	}
/*
	if n.Dir.IsDirectory() {
		return funcs
	} 
	return []func() {}
*/
}


func Expr(n Node, next Scanner) []func() {
	nt, _, arg := next(true)
	switch nt {
	case "(":
		return Expr(n, next)
	case "(file":
		return compose(file(n, Expr(n, next)), Expr(n, next))
	case "(dir":
		if !(n.Dir.IsDirectory()) {
			Expr(n, next)
			return []func() {}
		}
		return compose(dir(n, Expr(n, next)), Expr(n, next))
	case "(sub":
		subfuncs := []func() {}
		subs := subNodes(n)
		strs := strings.Fields(string(bytes))
		for i, _ := range subs {
			arg2 := arg
			tscan := func(use bool) (string, bool, int) {
				switch {
				case arg2 >= len(strs):
					return "", false, arg2
				case use:
					ret := strs[arg2]
					arg2++
					return ret, true, arg2
				}
				return strs[arg2], true, arg2
			}
			subfuncs = compose(sub(n, Expr(subs[i], tscan)), subfuncs)
		}
		Expr(n, next)
		return compose2(subfuncs, Expr(n, next))
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
	bytes, _ = ioutil.ReadFile("example2.ls")
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
		scanner := func(use bool) (string, bool, int) {
			switch {
			case arg >= len(strs):
				return "", false, arg
			case use:
				ret := strs[arg]
				arg++
				return ret, true, arg
			}
			return strs[arg], true, arg
		}

		result := Expr(node, scanner)
		for _, fn := range result {
			fn()
		}
	}
}

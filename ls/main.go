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
		Subs []Node
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
}


func Expr(n Node, next Scanner, norecurse bool) []func() {
	nt, _, arg := next(true)
	switch nt {
	case "(":
		return Expr(n, next, norecurse)
	case "(file":
		if !(n.Dir.IsRegular()) {
			Expr(n, next, true)
			return Expr(n, next, norecurse)
		}
		return compose(file(n, Expr(n, next, norecurse)), Expr(n, next, norecurse))
	case "(dir":
		if !(n.Dir.IsDirectory()) {
			Expr(n, next, true)
			//return Expr(n, next, norecurse)
			return []func() {}
		}
		return compose(dir(n, Expr(n, next, norecurse)), Expr(n, next, norecurse))
	case "(sub":
		subfuncs := []func() {}
		n.Subs = subNodes(n)
		strs := strings.Fields(string(bytes))
		for i, _ := range n.Subs {
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
			subfuncs = compose(sub(n, Expr(n.Subs[i], tscan, norecurse)), subfuncs)
		}
		Expr(n, next, true)
		return compose2(subfuncs, Expr(n, next, norecurse))
	case "(recurse)":
		if norecurse {
			return []func() {}
		}
		ret := compose2(doRecurse(n), Expr(n, next, norecurse))
		return ret
	case "(name)":
		return compose(name(n), Expr(n, next, norecurse))
	case "(nl)":
		return compose(nl(), Expr(n, next, norecurse))
	case ")":
		return []func(){}
	default:
		ftmp := func() {
			fmt.Printf(nt)
		}
		return compose(ftmp, Expr(n, next, norecurse))
	}
	return nil
}

func doRecurse(n Node) []func() {
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

	result := Expr(n, scanner, false)
	return result
/*
	for _, fn := range result {
		fn()
	}
*/
}

func main() {
	bytes, _ = ioutil.ReadFile("example2.ls")
	dname := os.Args[1]
	fpt, _ := os.Open(dname, os.O_RDONLY, 0666)
	names, _ := fpt.Readdirnames(-1)
	nodes := make([]Node, len(names))
	for i, cn := range names {
		nodes[i].Name = cn
		nodes[i].Dir, _ = os.Stat(cn)
	}

	node := new(Node)
	node.Name = dname
	node.Dir,_ = os.Stat(dname)

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

	result := Expr(*node, scanner, false)
	for _, fn := range result {
		fn()
	}

/*
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
*/
}

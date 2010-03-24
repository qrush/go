///////////////////////////////////////////////////////////////////////////////
// ls
// John Floren, Nick Quaranto
///////////////////////////////////////////////////////////////////////////////

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

	// Represents an element in the parseable .ls file
	Node struct {
		Name string
		Dir  *os.Dir
		Subs []Node
	}
)

var (
	// contents of the script file
	bytes []byte

	// precalculating byte sizes for os.Dir.Size checks
	a_kilobyte uint64 = 1024
	a_megabyte uint64 = a_kilobyte * 1024
	a_gigabyte uint64 = a_megabyte * 1024
)

///////////////////////////////////////////////////////////////////////////////
// Utility functions
///////////////////////////////////////////////////////////////////////////////

// Stick a function at the start of an array of functions
func compose(f func(), funcs []func()) []func() {
	newFuncs := make([]func(), len(funcs)+1)
	newFuncs[0] = f
	for i := range funcs {
		newFuncs[i+1] = funcs[i]
	}
	return newFuncs
}

// Join array of functions f1 in front of f2
func join(f1 []func(), f2 []func()) []func() {
	newFuncs := make([]func(), len(f1)+len(f2))
	for i := range f1 {
		newFuncs[i] = f1[i]
	}
	for i := range f2 {
		newFuncs[i+len(f1)] = f2[i]
	}
	return newFuncs
}

// print out error and die
func error(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

// Get the list of below the given node
func subNodes(n Node) ([]Node, os.Error) {
	if n.Dir.IsDirectory() {
		fpt, err := os.Open(n.Name, os.O_RDONLY, 0)
		if err != nil {
			return []Node{}, err
		}
		names, err2 := fpt.Readdirnames(-1)
		if err2 != nil {
			return []Node{}, err2
		}
		nodes := make([]Node, len(names))
		for i, cn := range names {
			nodes[i].Name = n.Name + "/" + cn
			nodes[i].Dir, err = os.Stat(nodes[i].Name)
			if err != nil {
				return []Node{}, err
			}
		}
		fpt.Close()
		return nodes, nil
	}
	return []Node{}, nil
}

// Recurse on a given node; this essentially starts the entire program over with a new node
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
}

///////////////////////////////////////////////////////////////////////////////
// Nodes
///////////////////////////////////////////////////////////////////////////////

// Print the name of the file
func name(n Node) func() { return func() { fmt.Printf("%s", n.Name) } }

// Return a function which prints out the given string
func print(s string) func() { return func() { fmt.Println(s) } }

// Return a function which prints out the given string and format
func printf(f string, a ...interface{}) func() {
	return func() { fmt.Printf(f, a) }
}

// Print the node's size in human-readable format
func humansize(n Node) func() {
	return func() {
		if n.Dir.Size < a_kilobyte {
			fmt.Printf("%v B", n.Dir.Size)
		} else if n.Dir.Size < a_megabyte {
			fmt.Printf("%2.0f KB", float64(n.Dir.Size)/float64(a_kilobyte))
		} else if n.Dir.Size < a_gigabyte {
			fmt.Printf("%2.1f MB", float64(n.Dir.Size)/float64(a_megabyte))
		} else {
			fmt.Printf("%2.1f GB", float64(n.Dir.Size)/float64(a_gigabyte))
		}
	}
}

// Returns a function for (file ... )
// If the node is a file, execute the list of functions
func file(n Node, funcs []func()) func() {
	return func() {
		if n.Dir.IsRegular() {
			for _, fn := range funcs {
				fn()
			}
		}
	}
}

// Returns a function for (dir ... )
// If node is a dir, execute the list of functions, otherwise do nothing
func dir(n Node, funcs []func()) func() {
	return func() {
		if n.Dir.IsDirectory() {
			for _, fn := range funcs {
				fn()
			}
		}
	}
}

// Execute (sub ... )
// Only works if node is a directory
func sub(n Node, funcs []func()) func() {
	return func() {
		if n.Dir.IsDirectory() {
			for _, fn := range funcs {
				fn()
			}
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// Expression Parser
// n is a file or directory
// next will give the next token
// norecurse should be set when the current node is a file and we see (dir ... ) or (sub ... )
// This flag is necessary because Expr is used to "eat" the remainder of the expression, but
// allowing recursion can create infinite loops
func Expr(n Node, next Scanner, norecurse bool) []func() {
	nt, ok, arg := next(true)
	if !ok {
		error("syntax error")
	}

	switch nt {
	case "(":
		return Expr(n, next, norecurse)
	case "(file":
		if !(n.Dir.IsRegular()) {
			// Not a file, so eat the rest of the expression without recursing
			Expr(n, next, true)
			return Expr(n, next, norecurse)
		}
		return compose(file(n, Expr(n, next, norecurse)), Expr(n, next, norecurse))
	case "(dir":
		if !(n.Dir.IsDirectory()) {
			// Not a directory, so eat the rest of the expression without recursing
			Expr(n, next, true)
			return []func(){}
		}
		return compose(dir(n, Expr(n, next, norecurse)), Expr(n, next, norecurse))
	case "(sub":
		subfuncs := []func(){}
		var err os.Error
		// Find the items under the current node
		n.Subs, err = subNodes(n)
		if err != nil {
			Expr(n, next, true)
			return compose(printf("Cannot get contents of %v, error = %v", n.Name, err), Expr(n, next, norecurse))
		}
		strs := strings.Fields(string(bytes))
		for i, _ := range n.Subs {
			// We create a new scanner because using "next" would 'use up' all the expressions on the first file.
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
		return join(subfuncs, Expr(n, next, norecurse))
	case "(recurse)":
		if norecurse {
			return []func(){}
		}
		ret := join(doRecurse(n), Expr(n, next, norecurse))
		return ret
	case "(tab)":
		return compose(print("	"), Expr(n, next, norecurse))
	case "(name)":
		return compose(name(n), Expr(n, next, norecurse))
	case "(user)":
		return compose(printf("%d", n.Dir.Uid), Expr(n, next, norecurse))
	case "(group)":
		return compose(printf("%d", n.Dir.Gid), Expr(n, next, norecurse))
	case "(size)":
		return compose(printf("%d", n.Dir.Size), Expr(n, next, norecurse))
	case "(human_size)":
		return compose(humansize(n), Expr(n, next, norecurse))
	case "(nl)":
		return compose(print(""), Expr(n, next, norecurse))
	case ")":
		return []func(){}
	default:
		// By default, we print all un-recognized items
		return compose(print(nt), Expr(n, next, norecurse))
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Run! Checks os.Args and gets the parser/evaluator started.
func main() {
	if len(os.Args) != 3 {
		error("Usage: ls [directory] [script.ls]")
	}

	dname := os.Args[1]
	if dname[len(dname)-1] == '/' {
		dname = dname[0 : len(dname)-1]
	}
	bytes, _ = ioutil.ReadFile(os.Args[2])

	var err os.Error
	node := new(Node)
	node.Name = dname
	node.Dir, err = os.Stat(dname)
	if err != nil {
		error("Couldn't stat " + dname)
	}

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
}

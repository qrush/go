/*

A Set is a DAG of targets which may have other targets as prerequisites.
This package implements set construction and depth-first traversal
of target trees.

A Target is constructed from one or more lines. The first line
contains the target name and optionally prerequisite names,
all separated by white space.

A Set is constructed from targets which are separated by one or more blank lines.

Note that prerequisites may overlap. The traversal will visit each prerequisite
only once.

*/
package dag

import (
	"io/ioutil"
	"flag"
	"fmt"
	"os"
	"strings"
)

type (
	// What a set of targets can do.
	Set interface {
		// Add or merge targets from the content of a file to the set;
		// on success return name of first target.
		AddFile(string, TargetFactory) (string, os.Error)

		// Add or merge targets from a string of lines to the set;
		// on success return name of first target.
		AddString(string, TargetFactory) (string, os.Error)

		// Add or merge targets from a list of lines to the set;
		// on success return name of first target.
		Add([]string, TargetFactory) (string, os.Error)

		// Add or merge one named target to the set;
		// on success return the target in the set.
		Put(Target) (Target, os.Error)

		// Return target corresponding to a name, if any.
		Get(string) Target

		// Send an action depth-first to all prerequisites
		// and then to a target itself.
		Apply(Target, Action) os.Error

		// Display.
		String() string
	}

	// What a target must do.
	Target interface {
		// Merge a target to a target (assumed of equal names);
		// on success return the receiver.
		Merge(Target) (Target, os.Error)

		// If not yet done, send an action to all prerequisites.
		ApplyPreq(Action) os.Error

		// If not yet done, set done and apply an action.
		Apply(Action) os.Error

		// Return name.
		Name() string

		// Display.
		String() string
	}

	// Create a target from a list of lines
	// and add or merges prerequisites, if any, to a set;
	// on success return the new target.
	TargetFactory func(Set, []string, TargetFactory) (Target, os.Error)

	// Manipulate a target.
	Action func(Target) os.Error

	Dag map[string]Target
)

// Contains the name of the dependency graph file.
var file string

// Bind file to the -f switch.
func init() {
	flag.StringVar(&file, "f", "mkfile", "file with target, sources, and command lines")
}

func (d Dag) AddFile(name string, fac TargetFactory) (string, os.Error) {
	fmt.Println("Adding a file: " + name)
	var bytes []byte
	first := ""

	file, err := os.Stat(name)

	if err == nil && file.IsRegular() {
		if bytes, err = ioutil.ReadFile(name); err == nil {
			blocks := strings.Split(string(bytes), "\n\n", 0)
			first, err = d.Add(blocks, fac)
		}
	}

	return first, err
}

func (d Dag) AddString(str string, fac TargetFactory) (string, os.Error) {
	fmt.Println("Adding strings: " + str)
	return "AddString", nil
}

func (d Dag) Add(strs []string, fac TargetFactory) (string, os.Error) {
	err := os.NewError("empty file")
	first := ""

	if len(strs) != 0 {
		for i, str := range strs {
			if i == 0 {
				first, err = d.AddString(str, fac)
			} else {
				_, err = d.AddString(str, fac)
			}
		}
	}

	return first, err
}

func (d Dag) Put(t Target) (Target, os.Error) { return nil, nil }

func (d Dag) Get(name string) Target { return nil }

func (d Dag) Apply(t Target, a Action) os.Error {
	return nil
}

func (d Dag) String() string { return "I'm a dag!" }

// Convenience method to run a typical command line.
// Must execute flag.Parse() before calling.
func Main(factory TargetFactory, action Action) {
	flag.Parse()
	s := new(Dag)
	if first, err := s.AddFile(file, factory); err != nil {
		fmt.Fprintf(os.Stderr, "dag: %s: %s\n", file, err)
		os.Exit(1)

	} else if flag.NArg() == 0 {
		fmt.Println("NArg() == 0")
		fmt.Println(first)
		s.Apply(s.Get(first), action)

	} else {
		fmt.Println("NArg() > 0")
		fmt.Println(flag.Args())
		for _, arg := range flag.Args() {
			if len(arg) > 0 {
				if target := s.Get(arg); target != nil {
					s.Apply(target, action)
				} else {
					fmt.Fprintf(os.Stderr, "dag: %s: undefined target\n", arg)
					os.Exit(1)
				}
			}
		}
	}
}

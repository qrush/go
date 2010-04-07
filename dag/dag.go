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
	"flag"
	"fmt"
	"os"
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
)

// Contains the name of the dependency graph file.
var file string

// Bind file to the -f switch.
func init() {
	flag.StringVar(&file, "f", "mkfile", "file with target, sources, and command lines")
}

// Convenience method to run a typical command line.
// Must execute flag.Parse() before calling.
func Main(factory TargetFactory, action Action) os.Error {
	s := NewSet()
	if first, err := s.AddFile(file, factory); err != nil {
		fmt.Fprintf(os.Stderr, "dag: %s: %s\n", file, err)
		os.Exit(1)

	} else if flag.NArg() == 0 {
		return s.Apply(s.Get(first), action)

	} else {
		for _, arg := range flag.Args() {
			if len(arg) > 0 {
				if target := s.Get(arg); target == nil {
					fmt.Fprintf(os.Stderr, "dag: %s: undefined target\n", arg)
					os.Exit(1)
				} else if err := s.Apply(target, action); err != nil {
				  return err
				}
			}
		}
	}
  return nil
}

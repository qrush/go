package dag

import (
	"io/ioutil"
	"fmt"
	"os"
	"strings"
)

type (
	Dag map[string]Target

	DagTarget struct {
		Field string	// Name of the target
		Prereqs map[string]*Target	// Prerequisites for the target
		Done bool	// True if target has been completed
	}
)

// Add or merge targets from the content of a file to the set;
// on success return name of first target.
func (d Dag) AddFile(name string, fac TargetFactory) (string, os.Error) {
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

// Add or merge targets from a string of lines to the set;
// on success return name of first target.
func (d Dag) AddString(str string, fac TargetFactory) (string, os.Error) {
	target, err := fac(d, strings.Split(str, "\n", 0), nil)

	if err == nil {
		return target.Name(), nil
	}

	return "", err
}

// Add or merge targets from a list of lines to the set;
// on success return name of first target.
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

// Add or merge one named target to the set;
// on success return the target in the set.
func (d Dag) Put(t Target) (Target, os.Error) {
	var err os.Error
	existing, ok := d[t.Name()]
	if ok {
		_, err = existing.Merge(t)
	} else {
		d[t.Name()] = t
	}
	return t, err
}

// Return target corresponding to a name, if any.
func (d Dag) Get(name string) Target {
	target, _ := d[name]
	return target
}

// Send an action depth-first to all prerequisites
// and then to a target itself.
func (d Dag) Apply(t Target, a Action) os.Error {
	return t.Apply(a)
}

// Return a string that displays the DAG object
func (d Dag) String() string {
	out := ""

	for key, value := range d {
		out = out + "KEY:\t" + key + "\t\tVALUE:\t" + value.String() + "\n"
	}

	return out
}

// Merge a DAG target into an existing target; return the receiver if successful
func (t *DagTarget) Merge(newt Target) (Target, os.Error) {
	dagTarget := (newt.(*DagTarget))

	if newt != nil {
		for _, prereq := range dagTarget.Prereqs {
			t.Prereqs[prereq.Name()] = prereq
		}
		return t, nil
	}
	return t, os.NewError("target parameter is nil")
}

// Apply the action to every prerequisite of the target.
func (t *DagTarget) ApplyPreq(a Action) os.Error {
	var err os.Error
	for _, prereq := range t.Prereqs {
		err = prereq.Apply(a)
		if err != nil {
			return err
		}
	}

	return err
}

// Apply the given action to the current target and its prereqs, if they exist.
func (t *DagTarget) Apply(a Action) os.Error {
	if ! t.Done {
		t.ApplyPreq(a)
		t.Done = true
		 return a(t)
	}

	return nil
}

// Return the name of the current target
func (t *DagTarget) Name() string {
	return t.Field
}

// Return a string showing the current target and all of its prerequisites
func (t *DagTarget) String() string {
	out := fmt.Sprintf("%s", t.Name())

	for _, prereq := range t.Prereqs {
		out = out + " " + prereq.Name() //fmt.Sprintf(" %s:%d", prereq.Name(), prereq.(*DagTarget).ID)
	}

	return out
}

// Make a brand new target with the given name.
func CreateDagTarget(field string) Target {
	var root Target
	root = new(DagTarget)
	root.(*DagTarget).Field = field
	root.(*DagTarget).Prereqs = make(map[string]*Target)

	return root
}

// Create a target from the given set of lines and add it to the specified Set.
func DagTargetFactory(s Set, lines []string, fac TargetFactory) (Target, os.Error) {
	var err os.Error
	fields := strings.Fields(lines[0])
	root := CreateDagTarget(fields[0])

	if len(fields) > 1 {
		for _, field := range fields[1:] {
			if et := s.Get(field); et != nil {
				root.(*DagTarget).Prereqs[field] = &et
			} else {
				tmp := CreateDagTarget(field)
				root.(*DagTarget).Prereqs[field] = &tmp
				_, err = s.Put(tmp)
			}
		}
	}

	_, err = s.Put(root)

	return root, err
}

// Return a new set of DAG targets.
func NewSet() Dag {
	return make(Dag)
}

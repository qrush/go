package dag

import (
	"io/ioutil"
	"os"
	"strings"
)

// Implements a Set.
type setImpl map[string]Target

// Constructor, returns an empty set.
func NewSet() Set { return make(setImpl) }

// Add or merge targets from the content of a file to the set;
// on success return name of first target.
func (this setImpl) AddFile(file string, factory TargetFactory) (string, os.Error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return this.AddString(string(bytes), factory)
}

// Add or merge targets from a string of lines to the set;
// on success return name of first target.
func (this setImpl) AddString(s string, factory TargetFactory) (string, os.Error) {
	return this.Add(strings.Split(s, "\n", 0), factory)
}

// Add or merge targets from a list of lines to the set;
// on success return name of first target.
func (this setImpl) Add(lines []string, factory TargetFactory) (string, os.Error) {
	begin := -1      // if >= 0, marks target/source line
	var first Target // if not nil

	// if there are lines add a target and maintain first
	create := func(end int) (Target, os.Error) {
		if begin >= 0 {
			if target, err := factory(this, lines[begin:end], factory); err != nil {
				return nil, err
			} else {
				if first == nil {
					first = target
				}
				begin = -1 // start over
				return this.Put(target)
			}
		}
		return nil, nil
	}

	for n, line := range lines {
		lines[n] = strings.TrimSpace(line)
		if len(lines[n]) == 0 { // found empty line
			if _, err := create(n); err != nil {
				return "", err
			}
		} else if begin < 0 { // found first non-empty line
			begin = n
		}
	}
	if _, err := create(len(lines)); err != nil {
		return "", err
	}
	return first.Name(), nil
}

// Add or merge one named target to the set;
// on success return the target in the set.
func (this setImpl) Put(target Target) (Target, os.Error) {
	if t, ok := this[target.Name()]; ok { // target is old
		return t.Merge(target)
	}
	this[target.Name()] = target
	return target, nil
}

// Return target corresponding to a name, if any.
func (this setImpl) Get(name string) Target {
	if target, ok := this[name]; ok {
		return target
	}
	return nil
}

// Send an action depth-first to all prerequisites
// and then to a target itself.
func (this setImpl) Apply(target Target, action Action) os.Error {
	// closure to recursively visit the prerequisites
	a := func(preq Target) os.Error { return this.Apply(preq, action) }
	if err := target.ApplyPreq(a); err != nil {
		return err
	}

	// visit the node
	return target.Apply(action)
}

// Display.
func (this setImpl) String() string {
	if len(this) == 0 {
		return ""
	}
	result := ""
	for _, t := range this {
		result += "\n\n" + t.String()
	}
	return result[2:]
}

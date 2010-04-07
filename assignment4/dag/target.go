package dag

import (
	"fmt"
	"os"
	"strings"
)

// Implements a Target.
type TargetImpl struct {
	name string            // target name
	Preq map[string]Target // prerequisites or nil
	mark bool              // used to detect cycles
	Done bool              // used to avoid overlaps
}

// An action to print a target.
func Print(target Target) os.Error {
	fmt.Println(target)
	return nil
}

// An action to reset a target so that it will be visited again.
func Reset(target Target) os.Error {
	target.(*TargetImpl).Done = false
	return nil
}

// Constructor: first line must contain one target and zero or more prerequisites,
// separated by white space. Unless already present, prerequisites are added to the set.
func NewTarget(set Set, lines []string, factory TargetFactory) (Target, os.Error) {
	result := new(TargetImpl)
	return result.Init(set, lines, factory)
}

// Initialization.
func (result *TargetImpl) Init(set Set, lines []string, factory TargetFactory) (Target, os.Error) {
	fields := strings.Fields(lines[0])

	// name, must exist
	result.name = fields[0]

	// Preq, if any
	if len(fields) > 1 {
		result.Preq = make(map[string]Target, len(fields)-1)
		for _, name := range fields[1:] {
			preq := set.Get(name)
			if preq == nil { // prerequisite will be new target
				if p, err := factory(set, []string{name}, factory); err != nil {
					return nil, err
				} else if preq, err = set.Put(p); err != nil {
					return nil, err
				}
			}
			result.Preq[name] = preq
		}
	}
	return result, nil
}

// Merge a target to a target (assumed of equal names);
// on success return the receiver.
func (this *TargetImpl) Merge(t Target) (Target, os.Error) {
	target := t.(*TargetImpl)

	// add or merge prerequisites (use target's values)
	if this.Preq == nil {
		this.Preq = target.Preq
	} else {
		for name, value := range target.Preq {
			this.Preq[name] = value
		}
	}

	return this, nil
}

// If not yet done, send an action to all prerequisites.
func (this *TargetImpl) ApplyPreq(action Action) os.Error {
	if !this.Done {
		if this.Preq != nil {
			if this.mark {
				return os.NewError(this.name + ": cyclic")
			}
			this.mark = true
			for _, p := range this.Preq {
				if err := action(p); err != nil {
					return err
				}
			}
			this.mark = false
		}
	}
	return nil
}

// If not yet done, set done and apply an action.
func (this *TargetImpl) Apply(action Action) os.Error {
	if !this.Done {
		this.Done = true
		return action(this)
	}
	return nil
}

// Return name.
func (this *TargetImpl) Name() string { return this.name }

// Display.
func (this *TargetImpl) String() string {
	result := this.name
	if this.Preq != nil {
		for name, _ := range this.Preq {
			result += " " + name
		}
	}
	return result
}

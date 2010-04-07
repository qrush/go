package mk

import (
	"dag"
	"fmt"
	"os"
)

/* Our new mk target will include a pointer to the old dag target */
type MkTarget struct {
	dag.Target
	Commands []string // The commands that accompany this target
}

// Create a new MkTarget.
func NewTarget(set dag.Set, lines []string, factory dag.TargetFactory) (dag.Target, os.Error) {
	result := new(MkTarget)
	return result.Init(set, lines, factory)
}

// Merge a target to a target (assumed of equal names);
// on success return the receiver.
func (this *MkTarget) Merge(t dag.Target) (dag.Target, os.Error) {
	target := t.(*MkTarget)
	tg := (target.Target).(*dag.TargetImpl)
	impl := (this.Target).(*dag.TargetImpl)

	// add or merge prerequisites (use target's values)
	if impl.Preq == nil {
		impl.Preq = tg.Preq
	} else {
		for name, value := range tg.Preq {
			impl.Preq[name] = value
		}
	}

	if len(this.Commands) == 0 {
		this.Commands = target.Commands
	}

	return this, nil
}

func (result *MkTarget) Init(set dag.Set, lines []string, factory dag.TargetFactory) (dag.Target, os.Error) {
	t, _ := dag.NewTarget(set, lines, factory)
	return &MkTarget{t, lines[1:]}, nil
}

func (this *MkTarget) Apply(action dag.Action) os.Error {
	t := this.Target.(*dag.TargetImpl)
	if !t.Done {
		t.Done = true
		return action(this)
	}
	return nil
}

func printcommands(c []string) {
	if len(c) != 0 {
		for _, s := range c {
			fmt.Println(s)
		}
	}
}

func Print(target dag.Target) os.Error {
	t := target.(*MkTarget)
	impl := t.Target.(*dag.TargetImpl)

	if dtop, err := os.Stat(impl.Name()); err == nil {
		for _, p := range impl.Preq {
			pt := p.(*MkTarget)
			pimpl := pt.Target.(*dag.TargetImpl)
			if d, err := os.Stat(pimpl.Name()); err == nil {
				if d.Mtime_ns > dtop.Mtime_ns {
					printcommands(t.Commands)
				}
			} else {
				printcommands(t.Commands)
			}
		}
	} else {
		printcommands(t.Commands)
	}
	return nil
}

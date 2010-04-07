package mk

import (
	"dag"
	"fmt"
	"os"
)

type MkTarget struct {
	dag.Target
	commands []string
}

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

	return this, nil
}

func (result *MkTarget) Init(set dag.Set, lines []string, factory dag.TargetFactory) (dag.Target, os.Error) {
	t, _ := dag.NewTarget(set, lines, factory)
	return &MkTarget{t, lines[1:]}, nil
}

func Print(target dag.Target) os.Error {
	fmt.Println("fffffuuuuuu")
	return nil
}

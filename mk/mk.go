package mk

import (
	"dag"
	"fmt"
	"os"
)

type MkTarget struct {
	dag.Target
	Commands []string
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

/*
func (this *MkTarget) ApplyPreq(action Action) os.Error {
	t := this.Target.(*dag.TargetImpl)
	if !t.Done {
		if t.Preq != nil {
			if t.mark {
				return os.NewError(t.name + ": cyclic")
			}
			t.mark = true
			for _, p := range t.Preq {
				if err := action(p); err != nil {
					return err
				}
			}
			t.mark = false
		}
	}
	return nil
}
*/

func (this *MkTarget) Apply(action dag.Action) os.Error {
	t := this.Target.(*dag.TargetImpl)
	if !t.Done {
		t.Done = true
		return action(this)
	}
	return nil
}

func Print(target dag.Target) os.Error {
	t := target.(*MkTarget)
	impl := t.Target.(*dag.TargetImpl)
	fmt.Println(impl.Name())
	fmt.Println(t.Commands)
	return nil
}

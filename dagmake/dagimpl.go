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
		Field string
		Prereqs map[string]*Target
		Done bool
	}
)

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

func (d Dag) AddString(str string, fac TargetFactory) (string, os.Error) {
	target, err := fac(d, strings.Split(str, "\n", 0), nil)

	if err == nil {
		return target.Name(), nil
	}

	return "", err
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

func (d Dag) Put(t Target) (Target, os.Error) {
	existing, ok := d[t.Name()]
	if ok {
		existing.Merge(t)
	} else {
		d[t.Name()] = t
	}
	return t, nil
}

func (d Dag) Get(name string) Target {
	target, _ := d[name]
	return target
}

func (d Dag) Apply(t Target, a Action) os.Error {
	return t.Apply(a)
}

func (d Dag) String() string {
	out := ""

	for key, value := range d {
		out = out + "KEY:\t" + key + "\t\tVALUE:\t" + value.String() + "\n"
	}

	return out
}

func (t *DagTarget) Merge(newt Target) (Target, os.Error) {
	dagTarget := (newt.(*DagTarget))

	for _, prereq := range dagTarget.Prereqs {
		t.Prereqs[prereq.Name()] = prereq
	}
	return t, nil
}

func (t *DagTarget) ApplyPreq(a Action) os.Error {
	for _, prereq := range t.Prereqs {
		prereq.Apply(a)
	}

	return nil
}

func (t *DagTarget) Apply(a Action) os.Error {
	if ! t.Done {
		t.ApplyPreq(a)
		t.Done = true
		a(t)
	}

	return nil
}

func (t *DagTarget) Name() string {
	return t.Field
}

func (t *DagTarget) String() string {
	out := fmt.Sprintf("%s", t.Name())

	for _, prereq := range t.Prereqs {
		out = out + " " + prereq.Name() //fmt.Sprintf(" %s:%d", prereq.Name(), prereq.(*DagTarget).ID)
	}

	return out
}

func CreateDagTarget(field string) Target {
	var root Target
	root = new(DagTarget)
	root.(*DagTarget).Field = field
	root.(*DagTarget).Prereqs = make(map[string]*Target)

	return root
}

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

func NewSet() Dag {
	return make(Dag)
}

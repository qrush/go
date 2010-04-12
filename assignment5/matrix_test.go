package matrix

import (
	"testing"
	"fmt"
)

func TestZeroMatrix(t *testing.T) {
	m, _ := Zeros(3, 3)
	if m.Rows() != 3 {
		t.Error("Wrong number of rows")
	}

	if m.Cols() != 3 {
		t.Error("Wrong number of columns")
	}

	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			if m.Get(i, j) != 0 {
				t.Error("Not zeroed matrix")
			}
		}
	}
}

func TestOnesMatrix(t *testing.T) {
	m, _ := Ones(2, 5)
	if m.Rows() != 2 {
		t.Error("Wrong number of rows")
	}

	if m.Cols() != 5 {
		t.Error("Wrong number of columns")
	}

	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			if m.Get(i, j) != 1 {
				t.Error("Not one'd matrix")
			}
		}
	}
}

func TestIHateNick(t *testing.T) {
        m, _ := Ones(4, 3)
        fmt.Println(m)
        n, _ := Ones(3, 3)
        fmt.Println(n)
        o, _ := Ones(3, 3)
        n.Add(o)
        fmt.Println(n)
        o.Add(n)
        fmt.Println(o)
        if added, err := m.Plus(o); err != nil {
                fmt.Println(err)
        } else {
                fmt.Println(added)
        }
        res, _ := n.Multiply(o)
        fmt.Println(res)

        zz,_ := Zeros(4,4)
        zz.Set(0,0,1)
        zz.Set(1,1,1)
        zz.Set(2,2,1)
        zz.Set(3,3,1)
        fmt.Println(zz)
        foo, _ := zz.Slice(0,2,0,2)
        fmt.Println(foo)
}

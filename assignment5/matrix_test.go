package matrix

import (
	"testing"
	"fmt"
	"strings"
)

func TestZeroMatrix(t *testing.T) {
	m, _ := Zeros(3, 3)
	if m.Rows() != 3 {
		t.Error("Wrong number of rows")
	}

	if m.Cols() != 3 {
		t.Error("Wrong number of columns")
	}

	checkAll(0, "Didn't zero properly", m, t)
}

func TestOnesMatrix(t *testing.T) {
	m, _ := Ones(2, 5)
	if m.Rows() != 2 {
		t.Error("Wrong number of rows")
	}

	if m.Cols() != 5 {
		t.Error("Wrong number of columns")
	}

	checkAll(1, "Didn't one properly", m, t)
}

func TestAddMatrix(t *testing.T) {
	m1, _ := Ones(3, 3)
	m2, _ := Ones(3, 3)
	m1.Add(m2)
	checkAll(2, "Didn't add to 2 properly", m1, t)
	m1.Add(m2)
	checkAll(3, "Didn't add to 3 properly", m1, t)
	m1.Add(m2)
	m1.Add(m2)
	checkAll(5, "Didn't add to 5 properly", m1, t)
}

func TestFourByThreeMatrixString(t *testing.T) {
	m, _ := Ones(4, 3)
	mstr := fmt.Sprintf("%s", *m)
	split := strings.Split(mstr[1:len(mstr)-1], "\n", 0)
	if len(split) != 4 {
		t.Error("Wrong number of lines")
	}

	for _, line := range split {
		str := strings.TrimSpace(line)
		if str != "[ 1.000000 1.000000 1.000000]" {
			t.Error("Wrong contents: '" + str + "'")
		}
	}
}

func TestThreeByThreeMatrixString(t *testing.T) {
	n, _ := Ones(3, 3)
	nstr := fmt.Sprintf("%s", *n)
	split := strings.Split(nstr[1:len(nstr)-1], "\n", 0)
	if len(split) != 3 {
		t.Error("Wrong number of lines")
	}

	for _, line := range split {
		str := strings.TrimSpace(line)
		if str != "[ 1.000000 1.000000 1.000000]" {
			t.Error("Wrong contents: '" + str + "'")
		}
	}
}

func TestAddingThreeByThreeMatrixString(t *testing.T) {
	n, _ := Ones(3, 3)
	o, _ := Ones(3, 3)
	n.Add(o)

	nstr := fmt.Sprintf("%s", *n)
	split := strings.Split(nstr[1:len(nstr)-1], "\n", 0)
	if len(split) != 3 {
		t.Error("Wrong number of lines")
	}

	for _, line := range split {
		str := strings.TrimSpace(line)
		if str != "[ 2.000000 2.000000 2.000000]" {
			t.Error("Wrong contents: '" + str + "'")
		}
	}
}

func TestAddingWrongSizes(t *testing.T) {
	n, _ := Ones(3, 3)
	o, _ := Ones(4, 3)
	n.Add(o)

	if _, err := n.Plus(o); err != nil {
		msg := fmt.Sprintf("%s", err)
		if msg != "Matrix dimensions do not match" {
			t.Error("Wrong error message: '" + msg + "'")
		}
	} else {
		t.Error("Did not return error for wrong matrix sizes")
	}
}

func checkAll(val float, error string, m *Matrix, t *testing.T) {
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			if m.Get(i, j) != val {
				t.Error(error)
			}
		}
	}
}

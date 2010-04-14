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

func TestMultiplyMatrix(t *testing.T) {
	m1, _ := Ones(3, 3)
	m1.Scale(4)
	m2, _ := Ones(3, 3)
	m2.Scale(3)
	m3, _ := Ones(3, 1)
	m3.Scale(2)

	m4, _ := m1.Multiply(m2)
	checkAll(36, "Didn't multiply to 36 properly", m4, t)

	m5, _ := m2.Multiply(m3)
	checkAll(18, "Didn't multiply to 18 properly", m5, t)
	if m5.Cols() != 1 {
		t.Error("Did not truncate columns")
	}
}

func TestBadMultiplyMatrix(t *testing.T) {
	m1, _ := Ones(3, 3)
	m2, _ := Ones(1, 2)
	_, err := m1.Multiply(m2)

	if err == nil {
		t.Error("Should complain about wrong row sizes")
	}
}

func TestPlusMatrix(t *testing.T) {
	m1, _ := Ones(3, 3)
	m2, _ := Ones(3, 3)
	m3, _ := m1.Plus(m2)

	checkAll(1, "Modified m1 matrix on plus", m1, t)
	checkAll(1, "Modified m2 matrix on plus", m1, t)
	checkAll(2, "Didn't add to 2 properly", m3, t)

	if &m1 == &m3 || &m2 == &m3 {
		t.Error("Copied reference instead of making new matrix")
	}
}

func TestScaleMatrix(t *testing.T) {
	m1, _ := Ones(3, 3)
	m1.Scale(4)
	checkAll(4, "Didn't scale to 4 properly", m1, t)
	m1.Scale(4)
	checkAll(16, "Didn't scale to 16 properly", m1, t)
}

func TestSliceMatrix(t *testing.T) {
	m1, _ := Ones(3, 3)

	firstrow, _ := m1.Slice(0, 1, 0, m1.Cols())
	otherrows, _ := m1.Slice(1, m1.Rows(), 0, m1.Cols())

	firstcol, _ := m1.Slice(0, m1.Rows(), 0, 1)
	othercols, _ := m1.Slice(0, m1.Rows(), 1, m1.Cols())

	// Modifying the original matrix just to make really sure
	m1.Scale(42)

	if firstrow.Rows() != 1 || firstrow.Cols() != 3 {
		t.Error("Should return one row with 3 columns")
	}
	checkAll(1, "Modified values during slice of first row", firstrow, t)

	if otherrows.Rows() != 2 || firstrow.Cols() != 3 {
		t.Error("Should return two rows with 3 columns")
	}
	checkAll(1, "Modified values during slice of other rows", otherrows, t)

	if firstcol.Rows() != 3 || firstcol.Cols() != 1 {
		t.Error("Should return three rows with 1 column")
	}
	checkAll(1, "Modified values during slice of first col", firstcol, t)

	if othercols.Rows() != 3 || othercols.Cols() != 2 {
		t.Error("Should return three rows with 2 columns")
	}
	checkAll(1, "Modified values during slice of other rows", othercols, t)
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

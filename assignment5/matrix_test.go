package matrix

import (
	"testing"
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

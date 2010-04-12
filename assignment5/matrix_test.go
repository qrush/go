package matrix

import (
	"testing"
)

func TestZeroMatrix(t *testing.T) {
	m, _ := Zeros(3, 3)
	if m.Rows() != 3 {
		t.Fail()
	}

	if m.Cols() != 3 {
		t.Fail()
	}
}

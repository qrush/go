package matrix

import (
	"os"
)

type MatrixInt interface {
	Add(MatrixInt) (MatrixInt, os.Error)
}

type Matrix struct {
	rows, cols int
	data       []float64
}

func Zeros(rows, cols int) (ret *Matrix, err os.Error) {
	ret = new(Matrix)
	ret.data = make([]float64, rows*cols)
	ret.rows = rows
	ret.cols = cols
	return
}

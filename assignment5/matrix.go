package matrix

import (
	"os"
	"fmt"
)


type MatrixInt interface {
	Add(*MatrixInt) os.Error // Modifies target matrix

	Plus(*MatrixInt) (*MatrixInt, os.Error) // Does not modify target

	Multiply(*MatrixInt) (*MatrixInt, os.Error) // Returns matrix product if successful

	Get(int, int) float // Get the element at the given row & column

	Set(int, int, float) // Set the specified element

	Rows() int

	Cols() int

	Slice(int, int, int, int) (*MatrixInt, os.Error)
}

type Matrix struct {
	rows, cols int
	data       []float
}

func Zeros(rows, cols int) (ret *Matrix, err os.Error) {
	ret = new(Matrix)
	ret.data = make([]float, rows*cols)
	ret.rows = rows
	ret.cols = cols
	return
}

func Ones(rows, cols int) (ret *Matrix, err os.Error) {
	ret = new(Matrix)
	ret.data = make([]float, rows*cols)
	for i := range ret.data {
		ret.data[i] = 1
	}
	ret.rows = rows
	ret.cols = cols
	return
}

func (this *Matrix) Rows() int { return this.rows }

func (this *Matrix) Cols() int { return this.cols }

func (this *Matrix) Add(m *Matrix) os.Error {
	if this.Rows() != m.Rows() || this.Cols() != m.Cols() {
		return os.NewError("Matrix dimensions do not match")
	}
	for i := 0; i < this.Rows(); i++ {
		for j := 0; j < this.Cols(); j++ {
			this.Set(i,j, this.Get(i,j) + m.Get(i,j))
		}
	}
	/*for i := range this.data {
		this.data[i] = this.data[i] + m.data[i]
	}*/
	return nil
}

func (this *Matrix) Plus(m *Matrix) (*Matrix, os.Error) {
	ret := new(Matrix)
	if this.rows != m.rows || this.cols != m.cols {
		return nil, os.NewError("Matrix dimensions do not match")
	}
	ret.data = make([]float, this.rows*this.cols)
	for i := range this.data {
		ret.data[i] = this.data[i] + m.data[i]
	}
	ret.rows = this.rows
	ret.cols = this.cols
	return ret, nil
}

func (this *Matrix) Get(row, col int) float {
	if (row < this.rows) && (col < this.cols) && (row >= 0) && (col >= 0) {
		return this.data[(row*this.cols)+col]
	}
	//return 0, os.NewError("Invalid row/col index")
	fmt.Println("invalid row/col index")
	return 0 // this is wrong
}

func (this *Matrix) Set(row, col int, val float) {
	this.data[(row*this.cols)+col] = val
}

func (this *Matrix) Multiply(m *Matrix) (*Matrix, os.Error) {
	if this.cols != m.rows {
		return nil, os.NewError("Invalid matrix dimensions for Multiply")
	}
	ret := new(Matrix)
	ret.data = make([]float, this.rows*m.cols)
	ret.rows = this.rows
	ret.cols = m.cols

	for i := 0; i < this.rows; i++ {
		for j := 0; j < m.cols; j++ {
			var sum float
			for k := 0; k < this.cols; k++ {
				sum += this.Get(i, k) * m.Get(k, j)
			}
			ret.Set(i, j, sum)
		}
	}
	return ret, nil
}

func (this *Matrix) Slice(rstart, rend, cstart, cend int) (*Matrix, os.Error) {
	if rstart >= rend || cstart >= cend {
		return nil, os.NewError("Invalid start/end specification")
	}
	ret, _ := Zeros(rend-rstart, cend-cstart)
	for i := 0; i < rend-rstart; i++ {
		for j := 0; j < cend-cstart; j++ {
			ret.Set(i, j, this.Get(i+rstart, j+cstart))
		}
	}
	return ret, nil
}

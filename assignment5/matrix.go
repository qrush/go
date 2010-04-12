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

func Zeros(rows, cols int) (*MatrixInt, os.Error) {
	var ret MatrixInt
	ret = new(Matrix)
	ret.(*Matrix).data = make([]float, rows*cols)
	ret.(*Matrix).rows = rows
	ret.(*Matrix).cols = cols
	return &ret, nil
}

func Ones(rows, cols int) (*MatrixInt, os.Error) {
	var ret MatrixInt
	ret = new(Matrix)
	ret.(*Matrix).data = make([]float, rows*cols)
	for i := range ret.(*Matrix).data {
		ret.(*Matrix).data[i] = 1
	}
	ret.(*Matrix).rows = rows
	ret.(*Matrix).cols = cols
	return &ret, nil
}

func (this *Matrix) String() string {
	var s string
	s += "[\n"
	for i := 0; i < this.Rows(); i++ {
		s += "	["
		for j := 0; j < this.Cols(); j++ {
			s += fmt.Sprintf(" %f", this.Get(i,j))
		}
		s += " ]\n"
	}
	s += "]\n"
	return s
}

func (this *Matrix) Rows() int { return this.rows }

func (this *Matrix) Cols() int { return this.cols }

func (this *Matrix) Add(m *MatrixInt) os.Error {
	if this.Rows() != m.Rows() || this.Cols() != m.Cols() {
		return os.NewError("Matrix dimensions do not match")
	}
	for i := 0; i < this.Rows(); i++ {
		for j := 0; j < this.Cols(); j++ {
			this.Set(i,j, this.Get(i,j) + m.Get(i,j))
		}
	}
	return nil
}

func (this *Matrix) Plus(m *MatrixInt) (*MatrixInt, os.Error) {
	if this.Rows() != m.Rows() || this.Cols() != m.Cols() {
		return nil, os.NewError("Matrix dimensions do not match")
	}
	//ret.data = make([]float, this.rows*this.cols)
	var ret *MatrixInt
	ret,_ = Zeros(this.Rows(), this.Cols())
	for i := 0; i < this.Rows(); i++ {
		for j := 0; j < this.Cols(); j++ {
			ret.Set(i,j, this.Get(i,j) + m.Get(i,j))
		}
	}
	return ret, nil
}

func (this *Matrix) Get(row, col int) float {
	if (row < this.Rows()) && (col < this.Cols()) && (row >= 0) && (col >= 0) {
		return this.data[(row*this.Cols())+col]
	}
	fmt.Println("invalid row/col index")
	return 0 // this is wrong
}

func (this *Matrix) Set(row, col int, val float) {
	this.data[(row*this.Cols())+col] = val
}

func (this *Matrix) Multiply(m *MatrixInt) (*MatrixInt, os.Error) {
	if this.Cols() != m.Rows() {
		return nil, os.NewError("Invalid matrix dimensions for Multiply")
	}
/*
	ret := new(Matrix)
	ret.data = make([]float, this.rows*m.cols)
	ret.rows = this.rows
	ret.cols = m.cols
*/
	ret,_ := Zeros(this.Rows(), this.Cols())

	for i := 0; i < this.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			var sum float
			for k := 0; k < this.Cols(); k++ {
				sum += this.Get(i, k) * m.Get(k, j)
			}
			ret.Set(i, j, sum)
		}
	}
	return ret, nil
}

func (this *Matrix) Slice(rstart, rend, cstart, cend int) (*MatrixInt, os.Error) {
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

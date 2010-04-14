// This package defines a Matrix interface and provides a dense matrix implementation.
package matrix

import (
	"os"
	"fmt"
)

type (
	// Define an interface for operations a Matrix must implement.
	Matrix interface {
		Add(*Matrix) os.Error // Modifies target matrix

		Plus(*Matrix) (*Matrix, os.Error) // Does not modify target

		Multiply(*Matrix) (*Matrix, os.Error) // Returns matrix product if successful

		Scale(float) // Scale the matrix by a given value.

		Get(int, int) float // Get the element at the given row & column

		Set(int, int, float) // Set the specified element

		Rows() int // Get the number of rows

		Cols() int // Get the number of columns

		Slice(int, int, int, int) (*Matrix, os.Error) // Get a slice of the specified matrix

		ForEach(Iterator) // Do an action on each cell of the matrix
	}

	// Implements a dense matrix
	DenseMatrix struct {
		rows, cols int
		data       []float
	}

	Iterator func(int, int)
)

// Return a new Matrix of all zeros
func Zeros(rows, cols int) (*Matrix, os.Error) {
	var ret Matrix
	ret = new(DenseMatrix)
	ret.(*DenseMatrix).data = make([]float, rows*cols)
	ret.(*DenseMatrix).rows = rows
	ret.(*DenseMatrix).cols = cols
	return &ret, nil
}

// Return a new Matrix of all ones
func Ones(rows, cols int) (*Matrix, os.Error) {
	ret, _ := Zeros(rows, cols)
	ret.ForEach(func(i, j int) { ret.Set(i, j, 1) })
	return ret, nil
}

// Return a string representing the matrix
func (this *DenseMatrix) String() string {
	s := "["
	for i := 0; i < this.Rows(); i++ {
		s += "	["
		for j := 0; j < this.Cols(); j++ {
			s += fmt.Sprintf(" %f", this.Get(i, j))
		}
		s += "]"
		if i != this.Rows()-1 {
			s += " \n"
		}
	}
	s += "	]"
	return s
}

// Run the given function on every cell of the matrix
func (this *DenseMatrix) ForEach(iter Iterator) {
	for i := 0; i < this.Rows(); i++ {
		for j := 0; j < this.Cols(); j++ {
			iter(i, j)
		}
	}
	return
}

// Get the number of rows in the matrix
func (this *DenseMatrix) Rows() int { return this.rows }

// Get the number of columns in the matrix
func (this *DenseMatrix) Cols() int { return this.cols }

// Add the argument matrix to the receiver.
func (this *DenseMatrix) Add(m *Matrix) os.Error {
	if this.Rows() != m.Rows() || this.Cols() != m.Cols() {
		return os.NewError("Matrix dimensions do not match")
	}

	this.ForEach(func(i, j int) { this.Set(i, j, this.Get(i, j)+m.Get(i, j)) })

	return nil
}

// Return a new Matrix containing the sum of the receiver and the argument.
func (this *DenseMatrix) Plus(m *Matrix) (*Matrix, os.Error) {
	if this.Rows() != m.Rows() || this.Cols() != m.Cols() {
		return nil, os.NewError("Matrix dimensions do not match")
	}
	var ret *Matrix
	ret, _ = Zeros(this.Rows(), this.Cols())
	this.ForEach(func(i, j int) { ret.Set(i, j, this.Get(i, j)+m.Get(i, j)) })
	return ret, nil
}

// Get the specified element.
func (this *DenseMatrix) Get(row, col int) float {
	if row < this.Rows() && col < this.Cols() && row >= 0 && col >= 0 {
		return this.data[(row*this.Cols())+col]
	}
	fmt.Println("invalid row/col index")
	return 0 // this is wrong
}

// Set the specified element.
func (this *DenseMatrix) Set(row, col int, val float) {
	this.data[(row*this.Cols())+col] = val
}

// Return the matrix consisting of the matrix product of the argument and the receiver.
func (this *DenseMatrix) Multiply(m *Matrix) (*Matrix, os.Error) {
	if this.Cols() != m.Rows() {
		return nil, os.NewError("Invalid matrix dimensions for Multiply")
	}
	ret, _ := Zeros(this.Rows(), m.Cols())

	ret.ForEach(func(i, j int) {
		var sum float
		for k := 0; k < this.Cols(); k++ {
			sum += this.Get(i, k) * m.Get(k, j)
		}
		ret.Set(i, j, sum)
	})
	return ret, nil
}

// Scale the receiver by the given value.
func (this *DenseMatrix) Scale(s float) {
	this.ForEach(func(i, j int) {
		t := this.Get(i, j)
		this.Set(i, j, t*s)
	})
}

// Return a new Matrix containing the specified range of the receiver. The starting rows and columns
// are inclusive, while the ending rows and columns are non-inclusive, as in the slicing operator. Thus,
// m.Slice(1, m.Rows(), 0, m.Cols()) to get rid of the first row.
func (this *DenseMatrix) Slice(rstart, rend, cstart, cend int) (*Matrix, os.Error) {
	if rstart >= rend || cstart >= cend {
		return nil, os.NewError("Invalid start/end specification")
	}
	ret, _ := Zeros(rend-rstart, cend-cstart)
	ret.ForEach(func(i, j int) { ret.Set(i, j, this.Get(i+rstart, j+cstart)) })
	return ret, nil
}

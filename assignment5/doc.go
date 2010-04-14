/*

The matrix package provides a simple interface to matrix operations.

Example usage:

// Creates a 3x3 matrix with all zeros
import "matrix"
m, err := matrix.Zeros(3, 3)

// Set the first member to 13
m.Set(0, 0, 13)

// Get the first member back (13)
m.Get(0, 0)

// Basic operations
m.Rows() // 3
m.Cols() // 3

// Print out the matrix
import "fmt"
fmt.Println(m)

[ [13.000000 0.000000 0.000000]
  [0.000000  0.000000 0.000000]
  [0.000000  0.000000 0.000000] ]

// Add another matrix to this matrix
n, _ := matrix.Ones()
m.Add(n)
fmt.Println(m)

[ [14.000000 1.000000 1.000000]
  [1.000000  1.000000 1.000000]
  [1.000000  1.000000 1.000000] ]

// Scale the matrix
m.Scale(2)

[ [28.000000 2.000000 2.000000]
  [2.000000  2.000000 2.000000]
  [2.000000  2.000000 2.000000] ]

// Get a slice of the matrix
// Arguments: starting row, ending row, starting col, ending col

m.Slice(0, m.Rows(), 0, m.Cols()) // Get the original matrix back
m.Slice(0, 1, 0, m.Cols())        // Just the first row

// Perform an action on each cell of the matrix
m.ForEach(func(i, j) {
  fmt.Println(m.Get(i, j))
})

// Multiply two matrices
o, _ := Ones(3, 3)
p, _ := Ones(3, 1)
o.Multiply(p)
fmt.Println(o)

[ [18.000000]
  [18.000000]
  [18.000000] ]

*/
package matrix

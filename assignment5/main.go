package main

import "matrix"
import "fmt"

func main() {
	m, _ := matrix.Ones(4, 3)
	fmt.Println(m)
	n, _ := matrix.Ones(3, 3)
	fmt.Println(n)
	o, _ := matrix.Ones(3, 3)
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

}

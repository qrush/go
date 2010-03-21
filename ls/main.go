package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type (
	// Scanner function, must return next token or false.
	// Advance past next token if argument is true.
	Scanner func(bool) (string, bool)

	// Parser function
	Parser func(Scanner) Eval

	// Interpreter function; returns value or os.Error.
	Eval func() interface{}
)

var strs string

func Expr(next Scanner) Eval {
	str, _ := next(true)
	fmt.Println(str)
	str, _ = next(true)
	fmt.Println(str)
	//if result := Or(next); result != nil {
	//  if _, ok := next(true); !ok { // eof
	//    return result
	//  }
	//}
	return nil
}

func main() {
	bytes, _ := ioutil.ReadFile("example1.ls")
	strs := strings.Fields(string(bytes))

	arg := 0 // consider os.Args[++arg] next

	result := Expr(func(use bool) (string, bool) {
		switch {
		case arg >= len(strs):
			return "", false
		case use:
			ret := strs[arg]
			arg++
			return ret, true
		}
		return strs[arg], true
	})

	//for _, str := range strs {
	fmt.Println(strs)
	fmt.Println(result)
	//}
}

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// The array of arguments, excluding the command name
var input []string

// Print an error and exit
func doerror(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

// Remove the current first argument from the list
func pop() {
	if len(input) > 1 {
		input = input[1:]
	}
}

// Find and evaluate an expression
func getExpr() float {
	t1 := getTerm()
	for input[0] == "+" || input[0] == "-" {
		op := input[0]
		pop()
		t2 := getTerm()
		if op == "+" {
			t1 += t2
		} else {
			t1 -= t2
		}
	}
	return t1
}

// Find an evaluate a term
func getTerm() float {
	f1 := getFactor()
	for input[0] == "*" || input[0] == "/" {
		op := input[0]
		pop()
		f2 := getFactor()
		if op == "*" {
			f1 *= f2
		} else if f2 != 0 {
			f1 /= f2
		} else {
			doerror("Divide by zero  FFFFFFFUUUUUUUUUuuuuuu......")
		}
	}
	return f1
}

// Find and evaluate a factor
func getFactor() float {

	if digitCheck, _ := regexp.MatchString("[0-9]+(\\.[0-9]+)?", input[0]); digitCheck {
		res, _ := strconv.Atof(input[0])
		pop()
		return res
	} else if input[0] == "(" {
		pop()
		res := getExpr()
		if input[0] != ")" {
			doerror("oh lawd no closing paren")
		}
		pop()
		return res
	}
	return 0
}

func main() {
	input = os.Args[1:]
	fmt.Println(getExpr())
}

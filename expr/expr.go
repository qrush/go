package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var input []string

func doerror(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func pop() {
	if len(input) > 1 {
		input = input[1:]
	}
}

func getExpr() int {
	t1 := getTerm()
	fmt.Println(t1)
	fmt.Println(input)
	for (input[0] == "+" || input[0] == "-") {
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

func getTerm() int {
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
			doerror("Divide by zero  FFFFFFFUUUUUUUUUuuuuuu......");
		}
	}
	return f1
}

func getFactor() int {

	if digitCheck,_ := regexp.MatchString("[0-9]+", input[0]); digitCheck {
		res,_ := strconv.Atoi(input[0])
		pop()
		return res
	} else if input[0] == "(" {
		pop()
		res := getExpr()
		if input[0] != ")" {
			doerror("oh lawd no closing paren")
		}
		return res
	}
	return 0
}	

func main() {
	input =  os.Args[1:]
	fmt.Println(getExpr())
}

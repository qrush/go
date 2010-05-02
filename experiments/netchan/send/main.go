package main

import "fmt"
import "netchan"
import "os"

func main() {
	fmt.Println("LOL")
	exp, eerr := netchan.NewExporter("tcp", ":9292")
	fmt.Println(eerr)
	c := make(chan os.SyscallError)
	err := exp.Export("octocator", c, netchan.Send, new(os.SyscallError))
	fmt.Println(err)
	c <- os.NewSyscallError("foo", 1)
}

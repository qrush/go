package main

import (
  "flag"
  "fmt"
)

var ip *int = flag.Int("flagname", 1234, "help message for flagname")

func main() {
  fmt.Printf("Testing out flags!\n");
  flag.Parse();
  fmt.Println("ip has value ", *ip);
}

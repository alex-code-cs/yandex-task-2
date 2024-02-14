package main

import (
	"fmt"
)

func main() {
	//	var expr = "2+             2 *  5"
	var expr = "(5+5) +* pi"
	var err = Parse(expr)
	fmt.Println(err)
}

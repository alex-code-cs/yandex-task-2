package main

import (
	"fmt"
)

func main() {
	var expr = "2+             2 *  5"
	var wrap = NewWrapper(expr)

	for wrap.Ch != CH_EOT {
		fmt.Print(string(wrap.Ch))
		wrap.NextChar()
	}

	fmt.Println()
}

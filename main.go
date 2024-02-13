package main

import (
	"fmt"
)

func main() {
	//	var expr = "2+             2 *  5"
	var expr = "(5+5) + pi"
	var lexer, err = NewLexer(expr)
	for lexer.Token != LEX_EOT {
		//	fmt.Println(lexer.Token)
		if err != nil {
			//fmt.Println(err)
		}
		err = lexer.NextLex()
	}
	err = Parse(expr)
	fmt.Println(err)
}

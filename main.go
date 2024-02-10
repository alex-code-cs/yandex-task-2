package main

import (
	"fmt"
)

func main() {
	//	var expr = "2+             2 *  5"
	var Lexer = NewLexer("pi")
	fmt.Println(Lexer.Name)
	fmt.Println(Lexer.Token)
}

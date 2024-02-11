package main

import (
	"fmt"
)

func main() {
	//	var expr = "2+             2 *  5"
	var lexer = NewLexer("pi *10 + 3")
	for lexer.Token != LEX_EOT {
		if lexer.Token == LEX_IDENT {
			fmt.Printf("%d\t%s\n", lexer.Token, lexer.Name)
		} else if lexer.Token == LEX_INT_NUMBER {
			fmt.Printf("%d\t%d\n", lexer.Token, lexer.IntValue.Int64())
		} else if lexer.Token == LEX_FLOAT_NUMBER {
			var floatValue, _ = lexer.FloatValue.Float64()
			fmt.Printf("%d\t%f\n", lexer.Token, floatValue)
		} else {
			fmt.Printf("%d\n", lexer.Token)
		}
		lexer.NextLex()
	}
}

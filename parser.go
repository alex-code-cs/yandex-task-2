package main

import (
	"fmt"
	"math/big"
)

const CH_EOT byte = 3

type Wrapper struct {
	Ch   byte
	data string
	pos  int
}

func (w *Wrapper) NextChar() {
	if w.pos == len(w.data) {
		w.Ch = CH_EOT
		return
	}
	w.Ch = w.data[w.pos]
	w.pos++
	for isWhitespace(w.Ch) && w.pos < len(w.data) {
		w.Ch = w.data[w.pos]
		w.pos++
	}
	if isWhitespace(w.Ch) {
		w.Ch = CH_EOT
	}
}

func NewWrapper(data string) *Wrapper {
	var wrap = Wrapper{data: data, pos: 0}
	wrap.NextChar()
	return &wrap
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t'
}

type Lex = int

const LEX_INT_NUMBER Lex = 0
const LEX_IDENT Lex = 1
const LEX_PLUS Lex = 2
const LEX_MINUS Lex = 3
const LEX_EQUAL Lex = 4
const LEX_MULTIPLY Lex = 5
const LEX_DIVIDE Lex = 6
const LEX_LBRACE Lex = 7
const LEX_RBRACE Lex = 8
const LEX_EOT Lex = 9
const LEX_NONE Lex = 10
const LEX_FLOAT_NUMBER = 11

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

type Lexer struct {
	Token      Lex
	Lex        rune
	IntValue   big.Int   // на случай если лексема - целое число
	FloatValue big.Float // на случай если лексема - дробное число
	Name       string    // на случай если лексема - идентификатор
	wrap       *Wrapper
}

// Не проверяет первый символ на букву
func (l *Lexer) ident() {
	l.Name = ""
	for isAlpha(l.wrap.Ch) {
		l.Name += string(l.wrap.Ch)
		l.wrap.NextChar()
	}
	l.Token = LEX_IDENT

}

// Не проверяет первый символ на цифру
func (l *Lexer) number() {

}

func (l *Lexer) NextLex() error {
	if isAlpha(l.wrap.Ch) {
		l.ident()
	} else if isNumber(l.wrap.Ch) {
		l.number()
	} else {
		switch l.wrap.Ch {
		case '+':
			l.Token = LEX_PLUS
		case '-':
			l.Token = LEX_MINUS
		case '*':
			l.Token = LEX_MULTIPLY
		case '/':
			l.Token = LEX_DIVIDE
		case '(':
			l.Token = LEX_LBRACE
		case ')':
			l.Token = LEX_RBRACE
		case '=':
			l.Token = LEX_EQUAL
		case CH_EOT:
			l.Token = LEX_EOT
		default:
			return fmt.Errorf("Неизвестный символ: %c", l.wrap.Ch)
		}
		l.wrap.NextChar()
	}
	return nil
}

func NewLexer(s string) *Lexer {
	var lexer = Lexer{}
	lexer.wrap = NewWrapper(s)
	lexer.NextLex()
	return &lexer
}

package main

import (
	"fmt"
	"math"
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

const LEX_NUMBER Lex = 0
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

func isNumber(c byte) bool {
	return c >= 48 && c <= 57
}

func isAlpha(c byte) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

type Lexer struct {
	Token      Lex
	Lex        rune
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
func (l *Lexer) number() error {
	var str string
	for isNumber(l.wrap.Ch) {
		str += string(l.wrap.Ch)
		l.wrap.NextChar()
	}
	if l.wrap.Ch == ',' {
		str += "."
		l.wrap.NextChar()
		if !(isNumber(l.wrap.Ch)) {
			if l.wrap.Ch == CH_EOT {
				return fmt.Errorf("Ожидалась цифра, но EOT")
			} else {
				return fmt.Errorf("Ожидалась цифра, но %c", l.wrap.Ch)
			}
		}
		for isNumber(l.wrap.Ch) {
			str += string(l.wrap.Ch)
			l.wrap.NextChar()
		}
	}
	l.Token = LEX_NUMBER
	l.FloatValue.SetString(str)
	return nil
}

func (l *Lexer) NextLex() error {
	if isAlpha(l.wrap.Ch) {
		l.ident()
	} else if isNumber(l.wrap.Ch) {
		err := l.number()
		if err != nil {
			return err
		}
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

func NewLexer(s string) (*Lexer, error) {
	var lexer = Lexer{}
	lexer.wrap = NewWrapper(s)
	var err = lexer.NextLex()
	if err != nil {
		return nil, err
	}
	return &lexer, nil
}

// #######       Парсер    ##############
/*

Выражение = Слагаемое {ОперСлож Слагаемое} (5+5)*5
Слагаемое = Множитель {ОперУмнож Множитель}
Множитель = Число | Идентификатор | "(" Выражение ")"
ОперСлож = "+" | "-"
ОперУмнож = "*" | "/"
Число = Цифра {Цифра}
Идентификатор = Буква {Буква}
*/

var lexer *Lexer
var nameTable map[string]float64

type Node struct {
	left      *Node
	right     *Node
	operation Lex
	value     big.Float
}

func SetNameTable() {
	nameTable = make(map[string]float64)
	nameTable["pi"] = math.Pi
}

func Parse(expr string) error {
	SetNameTable()
	var err error
	lexer, err = NewLexer(expr)
	if err != nil {
		return err
	}
	err = expression()
	if err != nil {
		return err
	}
	if lexer.Token != LEX_EOT {
		return fmt.Errorf("Ожидался конец текста")
	}

	return nil
}

// Выражение = Слагаемое {ОперСлож Слагаемое}
func expression(node **Node) error {
	var err = term(node) // слагаемое
	if err != nil {
		return err
	}
	for lexer.Token == LEX_MINUS || lexer.Token == LEX_PLUS {
		var tmp = node                  // запоминаем старое поддерево
		*node = &Node{}                 // создаем новую ноду
		(*node).operation = lexer.Token // не забываем про операцию
		(*node).left = *tmp             // подцепляем старое поддерево слева
		err = lexer.NextLex()
		if err != nil {
			return err
		}
		(*node) = &Node{}            // по соглашению создаем объект в памяти заранее
		err = term((&(*node).right)) // правую часть дерева отдаем на откуп term()
	}

	if err != nil {
		return err
	}

	return nil
}

// Слагаемое = Множитель {ОперУмнож Множитель}
func term(node **Node) error {
	var err = factor(node) //множитель
	if err != nil {
		return err
	}

	for lexer.Token == LEX_MULTIPLY || lexer.Token == LEX_DIVIDE {
		var tmp = node
		*node = &Node{}
		(*node).left = *tmp
		(*node).operation = lexer.Token
		err = lexer.NextLex()
		if err != nil {
			return err
		}
		(*node).right = &Node{}
		err = factor(&(*node).right)
		if err != nil {
			return err
		}
	}
	return nil
}

// Множитель = Число | Идентификатор | "(" Выражение ")"
func factor(node **Node) error {
	if lexer.Token == LEX_NUMBER {
		(*node).value = lexer.FloatValue // не забываем числовое значение
		(*node).operation = LEX_NONE
		var err = lexer.NextLex()
		if err != nil {
			return err
		}
	} else if lexer.Token == LEX_IDENT {
		var _, ok = nameTable[lexer.Name]
		if !ok {
			return fmt.Errorf("Неизвестный иденфикатор! %s", lexer.Name)
		}
		(*node).value.SetFloat64(nameTable[lexer.Name])
		var err = lexer.NextLex()
		if err != nil {
			return err
		}
	} else if lexer.Token == LEX_LBRACE {
		var err = lexer.NextLex()
		if err != nil {
			return err
		}
		err = expression(node)
		if err != nil {
			return err
		}
		if lexer.Token != LEX_RBRACE {
			return fmt.Errorf("Ожидалась скобка )")
		}
		lexer.NextLex()
	} else {
		return fmt.Errorf("Ожидалось число, имя или выражение в скобках")
	}
	return nil
}

func calculate(node *Node) big.Float {
	if node.left == nil && node.right == nil {
		return node.value
	}
	var leftResult big.Float = calculate(node.left)
	var rightResult big.Float = calculate(node.right)
	switch node.operation {
	case LEX_PLUS:
		return *(leftResult.Add(&leftResult, &rightResult))
	case LEX_MINUS:
		return *(leftResult.Sub(&leftResult, &rightResult))
	case LEX_DIVIDE:
		return *(leftResult.Quo(&leftResult, &rightResult))
	case LEX_MULTIPLY:
		return *(leftResult.Mul(&leftResult, &rightResult))
	}
}

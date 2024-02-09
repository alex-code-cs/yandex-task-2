package main

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

const LEX_NUMBER int = 0
const LEX_IDENT int = 1
const LEX_PLUS int = 2
const LEX_MINUS int = 3
const LEX_EQUAL int = 4
const LEX_MULTIPLY int = 5
const LEX_DIVIDE int = 6
const LEX_LBRACE int = 7
const LEX_RBRACE int = 8
const LEX_EOT int = 9
const LEX_NONE int = 10

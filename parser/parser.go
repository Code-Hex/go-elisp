package parser

import (
	"fmt"
	"strings"

	"github.com/Code-Hex/lisp/lexer"
	"github.com/Code-Hex/lisp/token"
)

type Printer interface {
	Print(indent int)
}

// Expression define to use type switch case
type Expression interface {
	Printer
}

// List makes list using car and cdr.
// ex. (a b c d)
// [a|]->[b|]->[c|]->[d|nil]
type List struct {
	Car Expression
	Cdr Expression
}

func (l *List) Print(indent int) {
	l.Car.Print(indent + 4)
	if l.Cdr != nil {
		l.Cdr.Print(indent + 4)
	} else {
		fmt.Println(strings.Repeat(" ", indent) + "nil")
	}
}

type ExpressionList []Expression

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(src string) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(src),
	}
}

func (p *Parser) Parse() (Expression, error) {
	tok, err := p.lexer.Lex()
	if err != nil {
		return nil, err
	}
	if tok.Type == token.EOF {
		return nil, nil
	}
	switch tok.Type {
	case token.LParen:
	case token.EOF:
	}
	return nil, nil
}

func (p *Parser) ParseList() (Expression, error) {
	return nil, nil
}

package parser

import (
	"fmt"
	"strings"

	"github.com/Code-Hex/go-elisp/lexer"
	"github.com/Code-Hex/go-elisp/token"
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

type Atom struct {
	Tok *token.Token
}

func (a *Atom) Print(indent int) {
	if a.Tok != nil {
		fmt.Println(strings.Repeat("-", indent) + a.Tok.String())
	}
}

type Integer struct {
	*Atom
	Value int
}

type Double struct {
	*Atom
	Value float64
}

type String struct {
	*Atom
	Value string
}

type Symbol struct {
	*Atom
	Value string
}

func (l *List) Print(indent int) {
	l.Car.Print(indent)
	if l.Cdr != nil {
		l.Cdr.Print(indent + 4)
	} else {
		fmt.Println(strings.Repeat("-", indent) + "nil")
	}
}

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(src string) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(src),
	}
}

func (p *Parser) ParseSExpr() (Expression, error) {
	sexp, err := p.ParseList()
	if err == nil {
		return sexp, nil
	}
	tok, err := p.lexer.Lex()
	if err != nil {
		return nil, err
	}
	if tok.Type == token.EOF {
		return nil, nil
	}
	switch tok.Type {
	case token.LParen:
		sexp, err := p.ParseList()
		if err != nil {
			return nil, err
		}
		return sexp, nil
	case token.EOF:
		return nil, nil
	}
	return nil, nil
}

func (p *Parser) ParseList() (Expression, error) {
	tok, err := p.lexer.Lex()
	if err != nil {
		return nil, err
	}
	switch tok.Type {
	case token.Symbol:
		cdr, err := p.ParseList()
		if err != nil {
			return nil, err
		}
		return &List{
			Car: &Symbol{
				Atom: &Atom{
					Tok: tok,
				},
				Value: tok.Literal,
			},
			Cdr: cdr,
		}, nil
	case token.RParen:
		return nil, nil
	default:
		panic(tok)
	}
}

package parser

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestParse(t *testing.T) {
	p := NewParser("(a b c)")
	sexp, err := p.ParseSExpr()
	if err != nil {
		t.Fatal(err)
	}
	sexp.Print(4)
	pp.Println(sexp)
}

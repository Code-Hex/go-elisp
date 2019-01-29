package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Code-Hex/lisp/token"
)

func TestLex(t *testing.T) {
	testcases := []struct {
		src      string
		expected []*token.Token
	}{
		{
			src: "(+ 4 5 1)",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "+",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "4",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "5",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "1",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
		{
			src: "(- -9 +2 3.00)",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "-",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "-9",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "+2",
				},
				&token.Token{
					Type:    "Double",
					Literal: "3.00",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
		{
			src: "(* 30.1234 #b1010111 3.2e+50)",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "*",
				},
				&token.Token{
					Type:    "Double",
					Literal: "30.1234",
				},
				&token.Token{
					Type:    "Binary",
					Literal: "#b1010111",
				},
				&token.Token{
					Type:    "Double",
					Literal: "3.2e+50",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
		{
			src: "(/ 1234 #o1234 #x1234)",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "/",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "1234",
				},
				&token.Token{
					Type:    "Oct",
					Literal: "#o1234",
				},
				&token.Token{
					Type:    "Hex",
					Literal: "#x1234",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
		{
			src: "(message \"hi\")",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "message",
				},
				&token.Token{
					Type:    "String",
					Literal: "\"hi\"",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
		{
			src: "(= (% n 2) 0)",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "=",
				},
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "%",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "n",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "2",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "0",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
		{
			src: "(= (% abcd [2 \"hello\" ?a]) ; comment\n0\n) ; this is comment",
			expected: []*token.Token{
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "=",
				},
				&token.Token{
					Type:    "LParen",
					Literal: "(",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "%",
				},
				&token.Token{
					Type:    "Symbol",
					Literal: "abcd",
				},
				&token.Token{
					Type:    "LBracket",
					Literal: "[",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "2",
				},
				&token.Token{
					Type:    "String",
					Literal: "\"hello\"",
				},
				&token.Token{
					Type:    "Char",
					Literal: "?a",
				},
				&token.Token{
					Type:    "RBracket",
					Literal: "]",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
				&token.Token{
					Type:    "Decimal",
					Literal: "0",
				},
				&token.Token{
					Type:    "RParen",
					Literal: ")",
				},
			},
		},
	}
	for _, c := range testcases {
		lexer := NewLexer(c.src)
		var tokens []*token.Token
		for {
			tok, err := lexer.Lex()
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if tok.Type == token.EOF {
				break
			}
			tokens = append(tokens, tok)
		}
		if diff := cmp.Diff(tokens, c.expected); diff != "" {
			t.Fatalf("(-got +want)\n%s\n", diff)
		}
	}
}

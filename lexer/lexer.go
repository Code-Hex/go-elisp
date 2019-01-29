package lexer

import (
	"bytes"
	"errors"
	"unicode"

	"github.com/Code-Hex/go-elisp/token"
)

type Lexer struct {
	src  []rune
	size int
	pos  int
	buf  bytes.Buffer
}

func NewLexer(src string) *Lexer {
	runes := []rune(src)
	return &Lexer{
		src:  runes,
		size: len(runes),
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= l.size {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) nextpeek() rune {
	pos := l.pos + 1
	if pos >= l.size {
		return 0
	}
	return l.src[pos]
}

func (l *Lexer) Lex() (*token.Token, error) {
	for {
		curCh := l.peek()
		switch curCh {
		case 0:
			return &token.Token{
				Type:    token.EOF,
				Literal: "EOF",
			}, nil
		case ';':
			l.skipComment()
		case '"':
			l.ScanString()
			str := l.Flush()
			return &token.Token{
				Type:    token.String,
				Literal: str,
			}, nil
		case '(':
			l.pos++
			return &token.Token{
				Type:    token.LParen,
				Literal: "(",
			}, nil
		case ')':
			l.pos++
			return &token.Token{
				Type:    token.RParen,
				Literal: ")",
			}, nil
		case '[':
			l.pos++
			return &token.Token{
				Type:    token.LBracket,
				Literal: "[",
			}, nil
		case ']':
			l.pos++
			return &token.Token{
				Type:    token.RBracket,
				Literal: "]",
			}, nil
		case '.':
			l.pos++
			return &token.Token{
				Type:    token.Dot,
				Literal: ".",
			}, nil
		case '\'':
			l.pos++
			return &token.Token{
				Type:    token.Quote,
				Literal: "'",
			}, nil
		case '?':
			return l.ScanChar(), nil
		case '+', '-':
			nextCh := l.nextpeek()
			if isDigit(nextCh) {
				l.buf.WriteRune(curCh)
				l.pos++
				return l.ScanNumbers(), nil
			}
			fallthrough
		case '*', '/', '%', '=':
			l.pos++
			return &token.Token{
				Type:    token.Symbol,
				Literal: string(curCh),
			}, nil
		default:
			switch {
			case unicode.IsSpace(curCh):
				l.skipBlanks() // white spaces or newlines, etc...
			case isDigit(curCh) || curCh == '#':
				return l.ScanNumbers(), nil
			case !isSymbol(curCh):
				return l.ScanSymbols(), nil
			default:
				return nil, errors.New("unreachable")
			}
		}
	}
}

func (l *Lexer) Flush() string {
	defer l.buf.Reset()
	return l.buf.String()
}

func (l *Lexer) skipBlanks() {
	for unicode.IsSpace(l.peek()) {
		l.pos++
	}
}

func (l *Lexer) skipComment() {
	for {
		ch := l.peek()
		if ch == '\n' || ch == 0 {
			break
		}
		l.pos++
	}
}

func (l *Lexer) ScanChar() *token.Token {
	// append head of '?'
	l.buf.WriteRune(l.peek())
	l.pos++
	if l.peek() == '\\' {
		l.buf.WriteRune(l.peek())
		l.pos++
	}
	l.buf.WriteRune(l.peek())
	l.pos++
	return &token.Token{
		Type:    token.Char,
		Literal: l.Flush(),
	}
}

func (l *Lexer) ScanString() {
	// append head of '"'
	l.buf.WriteRune(l.peek())
	l.pos++
	for {
		ch := l.peek()
		if ch == '"' {
			l.buf.WriteRune(l.peek())
			l.pos++
			break
		}
		if ch == '\\' {
			if l.nextpeek() == '"' {
				l.buf.WriteRune(l.peek())
				l.pos++
			}
		}
		l.buf.WriteRune(l.peek())
		l.pos++
	}
}

func (l *Lexer) ScanSymbols() *token.Token {
	for {
		ch := l.peek()
		if !(unicode.IsLetter(ch) || isDigit(ch)) {
			break
		}
		l.pos++
		l.buf.WriteRune(ch)
	}
	return &token.Token{
		Type:    token.Symbol,
		Literal: l.Flush(),
	}
}

func (l *Lexer) ScanNumbers() *token.Token {
	ch := l.peek()
	l.buf.WriteRune(ch)

	l.pos++
	// case of hex. e.g. 0x1234, 0xBadFace, etc...
	if ch == '#' && (l.peek() == 'x' || l.peek() == 'X') {
		l.buf.WriteRune(l.peek())
		l.pos++
		for {
			ch := l.peek()
			if !isHex(ch) {
				break
			}
			l.buf.WriteRune(ch)
			l.pos++
		}
		return &token.Token{
			Type:    token.Hex,
			Literal: l.Flush(),
		}
	}

	// case of octal. e.g. 0123, 0755, etc...
	if ch == '#' && (l.peek() == 'o' || l.peek() == 'O') {
		l.buf.WriteRune(l.peek())
		l.pos++
		for {
			ch := l.peek()
			if !isOctal(ch) {
				break
			}
			l.buf.WriteRune(ch)
			l.pos++
		}
		return &token.Token{
			Type:    token.Oct,
			Literal: l.Flush(),
		}
	}

	// case of octal. e.g. 0123, 0755, etc...
	if ch == '#' && (l.peek() == 'b' || l.peek() == 'B') {
		l.buf.WriteRune(l.peek())
		l.pos++
		for {
			ch := l.peek()
			if !isBinary(ch) {
				break
			}
			l.buf.WriteRune(ch)
			l.pos++
		}
		return &token.Token{
			Type:    token.Binary,
			Literal: l.Flush(),
		}
	}

	// case of integer, float. e.g. 10, 3020, 3.14, 1.4142, etc...
	{
		for {
			ch := l.peek()
			if !isDigit(ch) {
				break
			}
			l.buf.WriteRune(ch)
			l.pos++
		}

		// case of 32.45, 0.0038, 1.2e4, 3.2e+2, 2.4E-3
		if l.peek() == '.' {
			// append '.'
			l.buf.WriteRune(l.peek())
			l.pos++

			for {
				ch := l.peek()
				if !isDigit(ch) {
					break
				}
				l.buf.WriteRune(ch)
				l.pos++
			}
			// case of 1.2e4, 3.2e+2, 2.4E-3
			if l.peek() == 'e' || l.peek() == 'E' {
				l.buf.WriteRune(l.peek())
				l.pos++
				if l.peek() == '+' || l.peek() == '-' {
					l.buf.WriteRune(l.peek())
					l.pos++
				}
				for {
					ch := l.peek()
					if !isDigit(ch) {
						break
					}
					l.buf.WriteRune(ch)
					l.pos++
				}
			}
			return &token.Token{
				Type:    token.Double,
				Literal: l.Flush(),
			}
		}
	}

	return &token.Token{
		Type:    token.Decimal,
		Literal: l.Flush(),
	}
}

func isSymbol(ch rune) bool {
	return ch == '#' || ch == '\''
}
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isBinary(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isOctal(ch rune) bool {
	return '0' <= ch && ch <= '7'
}

func isHex(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

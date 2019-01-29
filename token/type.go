package token

type Type string

const (
	LParen   Type = "LParen"   // (
	RParen   Type = "RParen"   // )
	LBracket Type = "LBracket" // [ (vector type)
	RBracket Type = "RBracket" // ]
	Dot      Type = "Dot"      // .
	Quote    Type = "Quote"    // '
	String   Type = "String"   // "hello"
	Symbol   Type = "Symbol"   // let, char-to-string, \+1, 1+, \(*\ 1\ 2\), +-*/_~!@$%^&=:<>{}
	Decimal  Type = "Decimal"  // 1234
	Oct      Type = "Oct"      // 0123
	Hex      Type = "Hex"      // 0x1234
	Binary   Type = "Binary"   // 0b0101
	Double   Type = "Double"   // 10.123, 1.5e3
	Char     Type = "Char"     // ?a, ?\n
	EOF      Type = "EOF"
)

func (t Type) String() string {
	return string(t)
}

type Token struct {
	Type    Type
	Literal string
}

func (t Token) String() string {
	return t.Literal
}

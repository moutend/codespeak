package token

type Kind int

const (
	Invalid  Kind = 0
	Alphabet Kind = 1
	Number   Kind = 2
	Symbol   Kind = 3
	Unicode  Kind = 4
)

// GetKind returns kind of a rune.
func GetKind(r rune) Kind {
	if r < 0 {
		return Invalid
	}
	switch {
	case IsAlphabet(r):
		return Alphabet
	case IsNumber(r):
		return Number
	case IsSymbol(r):
		return Symbol
	}

	return Unicode
}

package token

// Kind represents a kind of token.
type Kind int

const (
	Invalid  Kind = 0
	Alphabet Kind = 1
	Number   Kind = 2
	Symbol   Kind = 3
	Unicode  Kind = 4
)

// IsAlphabet matches `[a-zA-Z]`.
func IsAlphabet(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

// IsNumber matches [0-9].
func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsSymbol matches non-alphanumeric characters.
func IsSymbol(r rune) bool {
	return r < 128 && !IsAlphabet(r) && !IsNumber(r)
}

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

// Token contains splitted words.
type Token struct {
	Kind Kind
	Text string
}

// Split converts string to Token slice.
func Split(s string) ([]Token, error) {
	if s == "" {
		return nil, nil
	}

	lastKind := Invalid
	tokens := []Token{}
	runes := []rune{}

	for _, r := range []rune(s) {
		if GetKind(r) != lastKind {
			tokens = append(tokens, Token{
				Kind: lastKind,
				Text: string(runes),
			})

			runes = []rune{}
		}

		runes = append(runes, r)
		lastKind = GetKind(r)
	}
	if len(runes) > 0 {
		tokens = append(tokens, Token{
			Kind: lastKind,
			Text: string(runes),
		})
	}

	return tokens[1:], nil
}

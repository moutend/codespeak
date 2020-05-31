package token

// Token contains splitted words.
type Token struct {
	Kind Kind
	Text string
}

// IsAlphabet matches [a-zA-Z].
func IsAlphabet(r rune) bool {
	if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
		return true
	}

	return false
}

// IsNumber matches [0-9].
func IsNumber(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}

	return false
}

// IsSymbol matches non-alphanumeric character.
func IsSymbol(r rune) bool {
	if r < 128 && !IsAlphabet(r) && !IsNumber(r) {
		return true
	}

	return false
}

// Split converts string to []token.Token data type.
func Split(s string) ([]Token, error) {
	if s == "" {
		return nil, nil
	}

	last := rune(-1)
	tokens := []Token{}
	runes := []rune{}

	for _, r := range []rune(s) {
		if GetKind(r) != GetKind(last) {
			tokens = append(tokens, Token{
				Kind: GetKind(last),
				Text: string(runes),
			})

			runes = []rune{}
		}

		runes = append(runes, r)
		last = r
	}
	if len(runes) > 0 {
		tokens = append(tokens, Token{
			Kind: GetKind(last),
			Text: string(runes),
		})
	}

	return tokens[1:], nil
}

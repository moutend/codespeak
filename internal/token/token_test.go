package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit1(t *testing.T) {
	actual, err := Split("")

	require.NoError(t, err)
	require.Nil(t, actual)
}

func TestSplit2(t *testing.T) {
	expected := []Token{
		{
			Kind: Symbol,
			Text: "@",
		},
	}

	actual, err := Split("@")

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSplit3(t *testing.T) {
	expected := []Token{
		{
			Kind: Alphabet,
			Text: "foo",
		},
		{
			Kind: Symbol,
			Text: " = ",
		},
		{
			Kind: Number,
			Text: "12345",
		},
		{
			Kind: Symbol,
			Text: "\n",
		},
	}

	actual, err := Split("foo = 12345\n")

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSplit4(t *testing.T) {
	expected := []Token{
		{
			Kind: Symbol,
			Text: "    ",
		},
		{
			Kind: Alphabet,
			Text: "foo",
		},
		{
			Kind: Symbol,
			Text: ".",
		},
		{
			Kind: Alphabet,
			Text: "Bar",
		},
		{
			Kind: Symbol,
			Text: "(",
		},
		{
			Kind: Number,
			Text: "123",
		},
		{
			Kind: Symbol,
			Text: ") // ",
		},
		{
			Kind: Unicode,
			Text: "これはコメントです",
		},
	}

	actual, err := Split("    foo.Bar(123) // これはコメントです")

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

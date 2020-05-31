package token

import (
	"testing"
)

func TestSplit1(t *testing.T) {
	actual, err := Split("")

	if err != nil {
		t.Fatal(err)
	}
	if actual != nil {
		t.Logf("actual: %+v\n", actual)
		t.Fatalf("wrong output")
	}
}

func TestSplit2(t *testing.T) {
	expected := []Token{
		{
			Kind: Symbol,
			Text: "@",
		},
	}
	actual, err := Split("@")

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(actual) {
		t.Logf("expected: %+v\n", expected)
		t.Logf("actual: %+v\n", actual)
		t.Fatalf("wrong output")
	}
	for i, v := range expected {
		if v.Kind != actual[i].Kind || v.Text != actual[i].Text {
			t.Logf("expected: %+v\n", v)
			t.Logf("actual: %+v\n", actual[i])
			t.Fatal("wrong output")
		}
	}
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

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(actual) {
		t.Logf("expected: %+v\n", expected)
		t.Logf("actual: %+v\n", actual)
		t.Fatalf("wrong output")
	}
	for i, v := range expected {
		if v.Kind != actual[i].Kind || v.Text != actual[i].Text {
			t.Logf("expected: %+v\n", v)
			t.Logf("actual: %+v\n", actual[i])
			t.Fatal("wrong output")
		}
	}
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

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(actual) {
		t.Logf("expected: %+v\n", expected)
		t.Logf("actual: %+v\n", actual)
		t.Fatalf("wrong output")
	}
	for i, v := range expected {
		if v.Kind != actual[i].Kind || v.Text != actual[i].Text {
			t.Logf("expected: %+v\n", v)
			t.Logf("actual: %+v\n", actual[i])
			t.Fatal("wrong output")
		}
	}
}

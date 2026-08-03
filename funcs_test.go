package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "mango      ",
			expected: []string{"mango"},
		},
		{
			input:    " mango      ",
			expected: []string{"mango"},
		},
	}

	for i, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("%v. different lengths %v != %v \n", i, len(actual), len(c.expected))
			t.Errorf("%v || %v", actual, c.expected)
			t.Fail()
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Fail()
				t.Errorf("words: %v || %v differ\n", word, expectedWord)
			}

		}
	}
}

package main 
import "testing"

func TestCleanInput(t *testing.T) {
    cases := []struct {
	input    string
	expected []string
    }{
    {
		input:    "  ",
		expected: []string{},
	},
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	{
		input:    "  Charmander Bulbasaur PIKACHU  ",
		expected: []string{"charmander", "bulbasaur", "pikachu"},
	},
    {
		input:    "  otro pokemon pokemon pokemon       world  ",
		expected: []string{"otro", "pokemon", "pokemon", "pokemon", "world"},
	},
    }

    for _, c := range cases {
	actual := cleanInput(c.input)
    if len(c.expected) != len(actual){
        t.Errorf("expected: %v, got: %v", len(c.expected), len(actual))
    }
	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
        if word != expectedWord {
            t.Errorf("expected: %v, got %v", expectedWord, word)
        }
	}
    }

}

package dice

import (
  "testing"
  "errors"
)

type testCase struct {
  name string
  input string
  expectedTokens []token
  expectedError error
}

var invalidTestCases = []testCase {
  {"Missing number after d in dice expression when operator follows", "1d+1", nil, errors.New("Dice expression was malformed for token 1d at position 1")},
  {"Missing number after d in dice expression when space follows", "1d + 1", nil, errors.New("Dice expression was malformed for token 1d at position 1")},
  {"Missing number after d in dice expression at end of input string", "1+2d", nil, errors.New("Dice expression was malformed for token 2d at position 3")},
  {"Multiple d/D in same expression", "1dd2+3", nil, errors.New("Multiple d/D in the same expression at poisiton 2")},
  {"Invalid characters in input string", "1c2+3", nil, errors.New("Unknown/invalid character (c) found at position 1")},
}


var validTestCases = []testCase {
  {"Only a number is a valid token", "10", []token{ { literal, "10" }, { eof, ""} }, nil},
  {"Only a dice expression that has the pattern XdY is a valid token", "1d6", []token{ { dice, "1d6" }, { eof, ""} }, nil},
  {"Only a dice expression that has the pattern dY is a valid token", "d6", []token{ { dice, "d6" }, { eof, ""} }, nil},
  {"Input has '(' and ')' around any number of terms", "(d6+1)*2", []token{ {operator, "("}, { dice, "d6" }, {operator, "+"}, {literal, "1"}, { operator, ")"}, {operator, "*"}, { literal, "2"}, { eof, ""} }, nil},
  {"Input has many '(' and ')' around any number of terms", "((d6+1)*2)+(2d12/2)", []token{ {operator, "("},{operator, "("}, { dice, "d6" }, {operator, "+"}, {literal, "1"}, { operator, ")"}, {operator, "*"}, { literal, "2"},{ operator, ")"}, { operator, "+"}, {operator, "("}, {dice, "2d12"}, {operator, "/"}, {literal, "2"}, {operator, ")"}, { eof, ""} }, nil},

  // below are valid input strings for the tokenize method but aren't valid in the lexer.
  {"Only an operator is a valid token", "-", []token{ { operator, "-" }, { eof, ""} }, nil},
  {"Only an operator is a valid token", "-*/", []token{ { operator, "-" }, { operator, "*" }, { operator, "/" }, { eof, ""} }, nil}, 
  {"Empty string produces only EOF token", "", []token{ { eof, ""} }, nil},
  {"Contains only valid literals, dice expressions, and operators in any order", "+1d4/", []token{ {operator, "+"}, {dice, "1d4"}, {operator, "/"}, {eof, ""}}, nil},

}


func TestTokenizeWithInvalidInputString(t *testing.T) {
  for _, tc := range invalidTestCases {
    t.Run(tc.name, func(t *testing.T) {
      actual, err := tokenize(tc.input)
      if actual != nil {
        t.Fatalf("Expected actual to be nil. Got %v\n", actual)
      }

      if err == nil {
        t.Fatalf("Expected err to not be nil but it was.\n")
      }
      t.Logf("error with message %s\n", err.Error())
    })
  }
}

func TestTokenizeWithValidInputString(t *testing.T) {
  for _, tc := range validTestCases {
    t.Run(tc.name, func(t *testing.T) {
      actual, err := tokenize(tc.input)
      if actual == nil {
        t.Fatalf("Expected actual to not be nil. error was %s\n", err.Error())
      }

      if err != nil {
        t.Fatalf("Expected err to be nil but err has message %s.\n", err.Error())
      }

      if len(actual) != len(tc.expectedTokens) {
        t.Fatalf("Expected to have %d tokens but tokenize returned %d\n", len(tc.expectedTokens), len(actual))
      }

      for idx, token := range actual {
        if token.value != tc.expectedTokens[idx].value {
          t.Fatalf("Tokens values aren't equal at index %d. Actual: \"%s\"  and expected: \"%s\"", idx, token.value, tc.expectedTokens[idx].value)
        }
      }
    })
  }
}

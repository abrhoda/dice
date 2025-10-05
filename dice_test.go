package dice

import (
	"errors"
	//"fmt"
	"testing"
)

func TestTokenizeWithValidInputString(t *testing.T) {
	for _, tc := range validTokenizeTestCases {
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

type parseTestCase struct {
	name                    string
	token                   token
	expectedValueLowerBound int
	expectedValueUpperBound int
	expcetedError           error
}

var invalidEvaluateTokenTestCases = []parseTestCase{
	{"Error when EOF token type has .evaluate() called", token{eof, ""}, 0, 0, errors.New("Token type 3 does not support evaluate.")},
	{"Error when operator token type has .evaluate() called", token{operator, "+"}, 0, 0, errors.New("Token type 2 does not support evaluate.")},
	{"Error when dice token type has no d/D and .evaluate() called", token{dice, "1"}, 0, 0, errors.New("Did not find 'd' or 'D' in token with type = dice and value = 1")},
	{"Error when dice token type has 2 d/D and .evaluate() called", token{dice, "1dD3"}, 0, 0, errors.New("strconv.Atoi(t.value[idx+1:]) should return error.")},
	{"Error when literal token type is a dice expression and .evalute() called", token{literal, "1d6"}, 0, 0, errors.New("strconv.Atoi(t.value) should return error.")},
}

var validEvaluateTokenTestCases = []parseTestCase{
	{"Int returned when literal token type with single digit value and .evaluate() called", token{literal, "6"}, 6, 6, nil},
	{"Int returned when literal token type with many digit value and .evaluate() called", token{literal, "100"}, 100, 100, nil},
	{"Int returned when dice token type has expression with no count and only faces.", token{dice, "d6"}, 1, 6, nil}, // 1-6
	{"Int returned when dice token type has expression with count and faces.", token{dice, "10d1"}, 10, 10, nil},
	{"Int returned when dice token type has expression with count, capital D, and faces.", token{dice, "2D4"}, 2, 8, nil},
}

func TestTokenEvaluateWithInvalidToken(t *testing.T) {
	for _, tc := range invalidEvaluateTokenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.token.evaluate()
			if err == nil {
				t.Fatalf("Expected err but was not. Actual value was %d.\n", actual)
			}
		})
	}
}

func TestTokenEvaluateWithValidToken(t *testing.T) {
	for _, tc := range validEvaluateTokenTestCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.token.evaluate()
			if err != nil {
				t.Fatalf("Expected err but was not. Actual value was %d. error: %s\n", actual, err.Error())
			}

			if actual == 0 {
				t.Fatalf("Actual was expected to have nonzero value.\n")
			}

			if actual < tc.expectedValueLowerBound {
				t.Fatalf("Actual value (%d) was less than lower bound (%d).\n", actual, tc.expectedValueLowerBound)
			}
			if tc.expectedValueUpperBound < actual {
				t.Fatalf("Actual value (%d) was more than upper bound (%d).\n", actual, tc.expectedValueUpperBound)
			}
		})
	}
}

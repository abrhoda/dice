package dice

import (
	"testing"
)

type parseTestCase struct {
	name           string
	input          []byte
	expectedTokens []token
	expectedResult int
}

var invalidParseTestCases = []parseTestCase{
	{"Malformed single term input returns error", []byte("("), []token{}, 0},
	{"Double operators in a row returns error", []byte("2**3"), []token{}, 0},
	{"Unmatched parens returns error", []byte("(1+1"), []token{}, 0},
	{"Invalid characters in input returns error", []byte("c+1"), []token{}, 0},
	{"Missing operator between terms in input returns error", []byte("1 1"), []token{}, 0},
}

var validParseTestCases = []parseTestCase{
	{"Single literal input returns value as int", []byte("1"), []token{{literal, "1"}, {eof, ""}}, 1},
	{"Single operator input returns value as int", []byte("1+3"), []token{{literal, "1"}, {operator, "+"}, {literal, "3"}, {eof, ""}}, 4},
	{"Multiple operator with same precedence input returns value as int", []byte("1+3-2"), []token{{literal, "1"}, {operator, "+"}, {literal, "3"}, {operator, "-"}, {literal, "2"}, {eof, ""}}, 2},
	{"Multiple operator with different precedence input returns value as int", []byte("12-3*2"), []token{{literal, "12"}, {operator, "-"}, {literal, "3"}, {operator, "*"}, {literal, "2"}, {eof, ""}}, 6},
	{"Multiple operator with different precedence and paren input returns value as int", []byte("(12-3)*2"), []token{{operator, "("}, {literal, "12"}, {operator, "-"}, {literal, "3"}, {operator, ")"}, {operator, "*"}, {literal, "2"}, {eof, ""}}, 18},
	// TODO test (2*6)-2*(2/3) and 12/(3+3)
}

// TODO implement these
var whiteSpaceParseTestCases = []parseTestCase{
	{"Empty input returns 0", []byte(""), []token{{eof, ""}}, 0},
	{"Blank input returns 0", []byte("       "), []token{{eof, ""}}, 0},
	{"Whitespace is stripped from input and returns value as int", []byte("\n1\v\r+   1 \t  "), []token{{eof, ""}}, 2},
}

func TestParseWithValidInputString(t *testing.T) {
	for _, tc := range validParseTestCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			res, err := p.Parse()
			if err != nil {
				t.Fatalf("Expected error to be nil but was present with message %s\n", err.Error())
			}

			if len(p.tokens) != len(tc.expectedTokens) {
				t.Fatalf("Length of parser's tokens (%d) and test case's expected token's (%d) don't match.\n", len(p.tokens), len(tc.expectedTokens))
			}

			for idx, tkn := range p.tokens {
				if tkn.kind != tc.expectedTokens[idx].kind {
					t.Fatalf("parse token's kind (%d) does not match test case's expected token's kind (%d) at index %d\n", tkn.kind, tc.expectedTokens[idx].kind, idx)
				}

				if tkn.value != tc.expectedTokens[idx].value {
					t.Fatalf("parse token's value (%s) does not match test case's expected token's kind (%s) at index %d\n", tkn.value, tc.expectedTokens[idx].value, idx)
				}
			}

			if res != tc.expectedResult {
				t.Fatalf("Result of %d does not match test case's expected result %d\n", res, tc.expectedResult)
			}
		})
	}
}

func TestParseWithInvalidInputString(t *testing.T) {
	for _, tc := range invalidParseTestCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("Expected error but found none.\n")
			}
		})
	}
}

package dice

import (
	"testing"
)

type parseTestCase struct {
	name                    string
	input                   []byte
	expectedTokens          []token
	expectedCurrentTokenPos int
	hasError                bool
}

var invalidParseTestCases = []parseTestCase{
	{"", []byte{}, []token{}, 0, true},
}

var validParseTestCases = []parseTestCase{}

func TestParseWithValidInputString(t *testing.T) {
	for _, tc := range validScannerTestCases {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}

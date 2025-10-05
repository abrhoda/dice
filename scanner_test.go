package dice

import (
	"errors"
	//	"fmt"
	"math/rand/v2"
	"testing"
)

type scannerTestCase struct {
	name           string
	input          string
	expectedTokens []token
	expectedError  error
}

var invalidScannerTestCases = []scannerTestCase{
	{"Missing number after d in dice expression when operator follows", "1d+1", nil, errors.New("Dice expression was malformed for token 1d at position 1")},
	{"Missing number after d in dice expression when space follows", "1d + 1", nil, errors.New("Dice expression was malformed for token 1d at position 1")},
	{"Missing number after d in dice expression at end of input string", "1+2d", nil, errors.New("Dice expression was malformed for token 2d at position 3")},
	{"Multiple d/D in same expression", "1dd2+3", nil, errors.New("Multiple d/D in the same expression at poisiton 2")},
	{"Invalid characters in input string", "1c2+3", nil, errors.New("Unknown/invalid character (c) found at position 1")},
}

var validScannerTestCases = []scannerTestCase{
	{"Only a number is a valid token", "10", []token{{literal, "10"}, {eof, ""}}, nil},
	{"Only a dice expression that has the pattern XdY is a valid token", "1d6", []token{{dice, "1d6"}, {eof, ""}}, nil},
	{"Only a dice expression that has the pattern dY is a valid token", "d6", []token{{dice, "d6"}, {eof, ""}}, nil},
	{"Input has '(' and ')' around any number of terms", "(d6+1)*2", []token{{operator, "("}, {dice, "d6"}, {operator, "+"}, {literal, "1"}, {operator, ")"}, {operator, "*"}, {literal, "2"}, {eof, ""}}, nil},
	{"Input has many '(' and ')' around any number of terms", "((d6+1)*2)+(2d12/2)", []token{{operator, "("}, {operator, "("}, {dice, "d6"}, {operator, "+"}, {literal, "1"}, {operator, ")"}, {operator, "*"}, {literal, "2"}, {operator, ")"}, {operator, "+"}, {operator, "("}, {dice, "2d12"}, {operator, "/"}, {literal, "2"}, {operator, ")"}, {eof, ""}}, nil},

	// below are valid input strings for the tokenize method but aren't valid in the lexer.
	{"Only an operator is a valid token", "-", []token{{operator, "-"}, {eof, ""}}, nil},
	{"Only an operator is a valid token", "-*/", []token{{operator, "-"}, {operator, "*"}, {operator, "/"}, {eof, ""}}, nil},
	{"Empty string produces only EOF token", "", []token{{eof, ""}}, nil},
	{"Contains only valid literals, dice expressions, and operators in any order", "+1d4/", []token{{operator, "+"}, {dice, "1d4"}, {operator, "/"}, {eof, ""}}, nil},

	// TODO
	//{"Converts 'D' in dice expression to lowercase when D is the first character", "D6", []token{{dice, "d6"}, {eof, ""}}, nil},
	//{"Converts 'D' in dice expression to lowercase when D in middle of token", "1D6", []token{{dice, "1d6"}, {eof, ""}}, nil},
}

func TestScannerWithInvalidInputString(t *testing.T) {
	for _, tc := range invalidScannerTestCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner{[]byte(tc.input), 0, 0}
			_, err := s.readToken()
			for err == nil {
				_, err = s.readToken()
			}
			if err == nil {
				t.Fatalf("Expected err to not be nil but it was.\n")
			}
		})
	}
}

func TestScannerWithValidInputString(t *testing.T) {
	for _, tc := range validScannerTestCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner{[]byte(tc.input), 0, 0}
			tkn, err := s.readToken()
			tokens := []token{tkn}

			for tkn.kind != eof && err == nil {
				tkn, err = s.readToken()
				tokens = append(tokens, tkn)
			}

			if err != nil {
				t.Fatalf("Expected no error but found error with message %s.\n", err.Error())
			}

			if tkn.kind != eof {
				t.Fatalf("Expected last token to be eof at the end of the expression.")
			}

			if len(tokens) != len(tc.expectedTokens) {
				t.Fatalf("Expected len(tokens) == len(tc.expectedTokens) but wasn't")
			}

			for idx, expectedToken := range tc.expectedTokens {
				if tokens[idx].kind != expectedToken.kind {
					t.Fatalf("Actual and Expected token at index %d have different types. Actual = %d while expected = %d\n", idx, tokens[idx].kind, expectedToken.kind)
				}
				if tokens[idx].value != expectedToken.value {
					t.Fatalf("Actual and Expected token at index %d have different values. Actual = %s while expected = %s\n", idx, tokens[idx].value, expectedToken.value)
				}
			}
		})
	}
}

type peekOrReadByteTestCase struct {
	name         string
	scannerInput string
	expectedByte byte
	expectedPos  int
}

var peekTestCases = []peekOrReadByteTestCase{
	{"Gives the only byte when input is 1 byte in length and does not advance the scanner's internal position", "1", '1', 0},
	{"Gives the first byte when input is more than 1 byte in length and does not advance the scanner's internal position", "d6", 'd', 0},
	{"Gives the EOF byte when input is 0 bytes in length and does not advance the scanner's internal position", "", EOF, 0},
}

var readTestCases = []peekOrReadByteTestCase{
	{"Gives the only byte when input is 1 byte in length and does advance the scanner's internal position", "1", '1', 1},
	{"Gives the first byte when input is more than 1 byte in length and not advance the scanner's internal position", "d6", 'd', 1},
	{"Gives the EOF byte when input is 0 bytes in length and does not advance the scanner's internal position", "", EOF, 0},
}

func TestScannerPeekByte(t *testing.T) {
	for _, tc := range peekTestCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner{[]byte(tc.scannerInput), 0, 0}
			p := s.peekByte()

			if p != tc.expectedByte {
				t.Fatalf("Expected to peek byte %c but peeked %c at index %d\n", tc.expectedByte, p, s.currentPos)
			}

			if s.currentPos != tc.expectedPos {
				t.Fatalf("Expected scanner's currentPos to be %d but was %d\n", tc.expectedPos, s.currentPos)
			}
		})
	}
}

func TestScannerReadByte(t *testing.T) {
	for _, tc := range peekTestCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner{[]byte(tc.scannerInput), 0, 0}
			p := s.peekByte()

			if p != tc.expectedByte {
				t.Fatalf("Expected to peek byte %c but peeked %c at index %d\n", tc.expectedByte, p, s.currentPos)
			}

			if s.currentPos != tc.expectedPos {
				t.Fatalf("Expected scanner's currentPos to be %d but was %d\n", tc.expectedPos, s.currentPos)
			}
		})
	}
}

type byteFunctionTestCase struct {
	name     string
	in       byte
	expected bool
}

var whitespaceBytes = []byte{' ', '\n', '\r', '\v', '\t'}
var digitBytes = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var diceCharacterBytes = []byte{'d', 'D'}
var operatorBytes = []byte{'+', '-', '*', '/', '(', ')'}

func getRandomByteOutsideSet(excludes []byte) byte {
	randomByte := byte(rand.IntN(128))
	found := false
	for {
		for _, b := range excludes {
			if b == randomByte {
				randomByte = byte(rand.IntN(128))
				found = true
				break
			}
		}
		if !found {
			return randomByte
		}
	}
}

var isWhitespaceTestCases = []byteFunctionTestCase{
	{"space is a whitespace char", whitespaceBytes[0], true},
	{"\\n (newline char) is a whitespace char", whitespaceBytes[1], true},
	{"\\r (return char) is a whitespace char", whitespaceBytes[2], true},
	{"\\v (vert tab char) is a whitespace char", whitespaceBytes[3], true},
	{"\\t (tab char) is a whitespace char", whitespaceBytes[4], true},
	{"Other char is not whitespace", getRandomByteOutsideSet(whitespaceBytes), false},
}

var isDigitTestCases = []byteFunctionTestCase{
	{"0-9 are valid digits", digitBytes[rand.IntN(len(digitBytes))], true},
	{"Chars outside 0-9 are not valid", getRandomByteOutsideSet(digitBytes), false},
}

var isDiceCharTestCases = []byteFunctionTestCase{
	{"d is a dice char", 'd', true},
	{"D is a dice char", 'D', true},
	{"Bytes that aren't d or D are not valid", getRandomByteOutsideSet(diceCharacterBytes), false},
}

var isOperatorTestCases = byteFunctionTestCase{}

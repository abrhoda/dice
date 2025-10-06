package dice

import (
	"testing"
)

type evaluateTestCase struct {
	name               string
	in                 token
	expectedLowerBound int
	expectedUpperBound int
}

var invalidEvaluateTestCases = []evaluateTestCase{
	{"Token with kind of eof returns error", token{eof, ""}, 0, 0},
	{"Token with kind of operator returns error", token{operator, "+"}, 0, 0},
	{"Token with kind of dice and without 'd/D' character returns error", token{dice, "6"}, 0, 0},
	// atoi actually fails this test because it splits on for d/D and then passes test to atoi as number
	{"Token with kind of dice and multiple 'd/D' characters returns error", token{dice, "1Dd6"}, 0, 0},
	{"Token with kind of dice and multiple 'd/D' characters and no 'count' prefix number returns error", token{dice, "Dd6"}, 0, 0},
	{"Token with kind of dice and multiple 'd/D' characters throughout the value returns error", token{dice, "1D2d6D"}, 0, 0},
	{"Token with kind of literal and non-digit characters returns error", token{dice, "61d11e"}, 0, 0},
}

func TestEvaluateWithInvalidToken(t *testing.T) {
	for _, tc := range invalidEvaluateTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.in.evaluate()
			if err == nil {
				t.Fatalf("Expected err to not be nil but it was.\n")
			}
		})
	}
}

var validEvaluateTestCases = []evaluateTestCase{
	{"Token with kind of literal single digit returns int value", token{literal, "1"}, 1, 1},
	{"Token with kind of literal multiple digits returns int value", token{literal, "1111"}, 1111, 1111},
	{"Token with kind of dice with value d{faces} returns int between 1 and {faces}", token{dice, "d4"}, 1, 4},
	{"Token with kind of dice with value D{faces} returns int between 1 and {faces}", token{dice, "D4"}, 1, 4},
	{"Token with kind of dice with value {count}*d{faces} returns int between {count} and {count}*{faces}", token{dice, "3d2"}, 3, 6},
	{"Token with kind of dice with value {count}*D{faces} returns int between {count} and {count}*{faces}", token{dice, "3d2"}, 3, 6},
}

func TestEvaluateWithValidToken(t *testing.T) {
	for _, tc := range validEvaluateTestCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.in.evaluate()
			if err != nil {
				t.Fatalf("Expected err to be nil but err had message %s.\n", err.Error())
			}

			if res < tc.expectedLowerBound {
				t.Fatalf("Result (%d) was less than expectedLowerBound (%d)\n", res, tc.expectedLowerBound)
			}

			if res > tc.expectedUpperBound {
				t.Fatalf("Result (%d) was more than expectedUpperBound (%d)\n", res, tc.expectedUpperBound)
			}
		})
	}
}

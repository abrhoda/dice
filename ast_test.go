package dice

import "testing"

type walkTestCase struct {
	name           string
	root           node
	expectedResult int
}

var invalidWalkTestCases = []walkTestCase{
	{"Operator token without left returns an error", node{token{operator, "+"}, nil, &node{token{literal, "1"}, nil, nil}}, 0},
	{"Operator token without right returns an error", node{token{operator, "+"}, &node{token{literal, "1"}, nil, nil}, nil}, 0},
	{"Recursively, when node is missing left returns an error", node{token{literal, "+"}, &node{token{literal, "1"}, &node{token{literal, "1"}, nil, nil}, nil}, &node{token{literal, "1"}, nil, nil}}, 0},
	{"Recursively, when node is missing right returns an error", node{token{literal, "+"}, &node{token{literal, "1"}, nil, nil}, &node{token{literal, "1"}, &node{token{literal, "1"}, nil, nil}, nil}}, 0},
}

var validWalkTestCases = []walkTestCase{}

func TestWalkWithValidAst(t *testing.T) {
	for _, tc := range validWalkTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := walk(&tc.root)
			if err != nil {
				t.Fatalf("Expected error to be nil but got error with message %s\n", err.Error())
			}
		})
	}
}

func TestWalkWithInvalidAst(t *testing.T) {
	for _, tc := range invalidWalkTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := walk(&tc.root)
			if err == nil {
				t.Fatalf("Expected error to not be nil\n")
			}
		})
	}
}

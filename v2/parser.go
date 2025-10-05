package main

// binding power type and mapping
type weight struct {
	lbp float64
	rbp float64
}

var operatorBindingPowers = map[string]weight{
	"+": {1.0, 1.1},
	"-": {1.0, 1.1},
	"*": {2.0, 2.1},
	"/": {2.0, 2.1},
}

// Picking a defualt slice size that will fit most common dice expressions
const DefaultTokenSliceSize = 10

type Parser struct {
	buffer         []byte
	tokens         []token
	curentTokenPos int
	//isError        bool
}

func NewParser(buffer []byte) Parser {
	return Parser{buffer, make([]token, 0, DefaultTokenSliceSize), 0}
}

// NOTE assemble in ast. Should dice tokens just be evaluated when putting them into the ast?

func (parser Parser) Parse() (int, error) {

	// scan tokens into parser.tokens
	s := newScanner(parser.buffer)
	for {
		t, err := s.readToken()
		if err != nil {
			return 0, err
		}
		parser.tokens = append(parser.tokens, t)

		if t.tType == eof {
			break
		}
	}

	return 0, nil
}

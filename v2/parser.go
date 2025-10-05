package main

import "fmt"

// Picking a defualt slice size that will fit most common dice expressions
const DefaultTokenSliceSize = 10

type Parser struct {
	buffer          []byte
	tokens          []token
	currentTokenPos int
	//isError        bool
}

func NewParser(buffer []byte) Parser {
	return Parser{buffer, make([]token, 0, DefaultTokenSliceSize), 0}
}

// NOTE assemble in ast. Should dice tokens just be evaluated when putting them into the ast?
func (parser *Parser) Parse() error {
	// scan tokens into parser.tokens
	s := newScanner(parser.buffer)
	for {
		t, err := s.readToken()
		if err != nil {
			return err
		}
		parser.tokens = append(parser.tokens, t)

		if t.kind == eof {
			break
		}
	}

	for i, t := range parser.tokens {
		fmt.Printf("%d. %v\n", i, t)
	}

	return nil
}

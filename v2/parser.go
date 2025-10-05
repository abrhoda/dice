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

// node type and functions
type node struct {
	token token
	left  *node
	right *node
}

type Parser struct {
	buffer      []byte
	curTokenPos int
}

func NewParser(buffer []byte) Parser {
	return Parser{buffer, 0}
}

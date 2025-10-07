package dice

import "fmt"

// Picking a defualt slice size that will fit most common dice expressions
const defaultTokenSliceSize = 10

// TODO make `.buffer` into `.Buffer` to allow switching out the byte buffer instead of needing to create new Parser structs to parse mutliple inputs.
// would also need to provide a `.Reset` function to allow the internal state of the Parser (mainly the tokens and currentTokenPos) to be zeroed out.
// buffer could not be reset/zeroed because the intention would be to swap that out after calling .Reset()
type parser struct {
	buffer          []byte
	tokens          []token
	currentTokenPos int
}

func NewParser(buffer []byte) parser {
	return parser{buffer, make([]token, 0, defaultTokenSliceSize), 0}
}



type weight struct {
	left  float64
	right float64
}

var operatorWeights = map[string]weight{
	"+": {1.0, 1.1},
	"-": {1.0, 1.1},
	"*": {2.0, 2.1},
	"/": {2.0, 2.1},
}

func (p *parser) astFromTokens(mbp float64) (*node, error) {
	var err error
	root := &node{p.tokens[p.currentTokenPos], nil, nil}
	p.currentTokenPos++
	if root.token.kind == operator && root.token.value == "(" {
		root, err = p.astFromTokens(0.0)
		if err != nil {
			return nil, err
		}
		temp := p.tokens[p.currentTokenPos]
		p.currentTokenPos++
		if temp.value != ")" {
			return nil, fmt.Errorf("Expression should have closing paren but none were found.")
		}

	} else if root.token.kind != dice && root.token.kind != literal {
		return nil, fmt.Errorf("Expression must start with a dice or literal. Found %d.", root.token.kind)
	}

	for {
		if p.currentTokenPos == len(p.tokens) {
			return root, nil
		}

		eofOrOp := p.tokens[p.currentTokenPos]
		if eofOrOp.kind == eof || eofOrOp.value == ")" {
			return root, nil
		} else if eofOrOp.kind == dice || eofOrOp.kind == literal {
			return nil, fmt.Errorf("Expected EOF or operation token. Found %d with value %s", eofOrOp.kind, eofOrOp.value)
		}

		lbp, rbp := operatorWeights[eofOrOp.value].left, operatorWeights[eofOrOp.value].right
		if lbp < mbp {
			break
		}

		p.currentTokenPos++
		rhs, err := p.astFromTokens(rbp)
		if err != nil {
			return nil, err
		}

		root = &node{eofOrOp, root, rhs}
	}

	return root, nil
}

func (parser *parser) Parse() (int, error) {
	// TODO when making the `.Reset` func for zeroing out, check that we have parser.currentTokenPos == 0 here or return an error
	// Maybe have a bool autoCleanUp param and if set to true, defer a call to a .reset func to the end of this func automatically
	// if parser.currentTokenPos != 0 {
	// 	return 0, fmt.Errorf("Detected reuse of parser without calling .Reset().")
	// }
	s := scanner{parser.buffer, 0, 0}
	for {
		t, err := s.readToken()
		if err != nil {
			return 0, err
		}
		parser.tokens = append(parser.tokens, t)

		if t.kind == eof {
			break
		}
	}

	ast, err := parser.astFromTokens(0.0)
	if err != nil {
		return 0, err
	}

	return walk(ast)
}

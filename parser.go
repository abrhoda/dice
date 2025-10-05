package dice

import "fmt"

// Picking a defualt slice size that will fit most common dice expressions
const DefaultTokenSliceSize = 10

type Parser struct {
	Buffer          []byte
	tokens          []token
	currentTokenPos int
	//isError        bool
}

func NewParser(buffer []byte) Parser {
	return Parser{buffer, make([]token, 0, DefaultTokenSliceSize), 0}
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

func (p *Parser) astFromTokens(mbp float64) (*node, error) {
	var err error
	root := &node{p.tokens[p.currentTokenPos], nil, nil}
	fmt.Printf("at start of astFromTokens = %v\n", root.token)
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

		fmt.Printf("eofOrOp = %v\n", eofOrOp)
		p.currentTokenPos++
		rhs, err := p.astFromTokens(rbp)
		if err != nil {
			return nil, err
		}

		root = &node{eofOrOp, root, rhs}
	}

	return root, nil
}

func (parser *Parser) Parse() (int, error) {
	s := scanner{parser.Buffer, 0, 0}
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

package main

import "fmt"

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

type node struct {
	token token
	left  *node
	right *node
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

func walk(root *node) (int, error) {
	if root.token.kind == eof {
		return 0, nil
	}
	fmt.Printf("root node: %v\n", root)

	if root.token.kind == operator {
		if root.left == nil || root.right == nil {
			return 0, fmt.Errorf("root node is an operator node with a nil right or left.")
		}
		lhs, err := walk(root.left)
		if err != nil {
			return 0, err
		}
		rhs, err := walk(root.right)
		if err != nil {
			return 0, err
		}

		switch root.token.value {
		case "+":
			return lhs + rhs, nil
		case "-":
			return lhs - rhs, nil
		case "*":
			return lhs * rhs, nil
		case "/":
			return lhs / rhs, nil
		default:
			return 0, fmt.Errorf("Invalid operator value found for token. Value was %s but should be +, -, *, or /.", root.token.value)
		}
	}
	return root.token.evaluate()
}

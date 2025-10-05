package main

//func iterNodes(nodes []node) {
//	for idx, node := range nodes {
//		fmt.Printf("%d - %d - %s, ", idx, node.tokenType, node.value)
//	}
//	fmt.Println("")
//}
//
//func printNode(tag string, node node) {
//	fmt.Printf("%s: type = %d, value = %s\n", tag, node.tokenType, node.value)
//}
//
//func buildAstFromNodes(nodes []node, mbp float64) (*node, error) {
//	var err error
//	root, rest := &nodes[0], nodes[1:]
//
//	printNode("root at start", *root)
//	if root.value == "(" {
//		//fmt.Printf("calling buildast without node with value %s\n", root.value)
//		root, err = buildAstFromNodes(rest, 0.0)
//		if err != nil {
//			return nil, err
//		}
//		var closing *node
//		closing, rest = &rest[0], rest[1:]
//		if closing.value != ")" {
//			return nil, fmt.Errorf("No closing paren was found. Found %s instead.", closing.value)
//		}
//	}
//	if root.tokenType != dice && root.tokenType != literal {
//		return nil, fmt.Errorf("Expression must start with a dice or literal. Found %d with value %s.", root.tokenType, root.value)
//	}
//
//	for {
//		if len(rest) == 0 {
//			return root, nil
//		}
//		printNode("start of for", *root)
//		eofOrOp := rest[0]
//		printNode("eofOrOp", eofOrOp)
//		if eofOrOp.tokenType == eof || eofOrOp.value == ")" {
//			return root, nil
//		} else if eofOrOp.tokenType == dice || eofOrOp.tokenType == literal {
//			return nil, fmt.Errorf("Expected EOF or operation token. Found %d with value %s", eofOrOp.tokenType, eofOrOp.value)
//		}
//
//		lbp, rbp := operatorBindingPowers[eofOrOp.value].lbp, operatorBindingPowers[eofOrOp.value].rbp
//		if lbp < mbp {
//			break
//		}
//
//		rest = rest[1:]
//		rhs, err := buildAstFromNodes(rest, rbp)
//
//		if err != nil {
//			return nil, err
//		}
//
//		root = &node{eofOrOp.tokenType, eofOrOp.value, root, rhs}
//	}
//	return root, nil
//}
//
//func evaluateAst(root *node) (int, error) {
//	if root.tokenType == eof {
//		return 0, nil
//	}
//
//	if root.tokenType == operator {
//		if root.left == nil || root.right == nil {
//			return 0, fmt.Errorf("root node is an operator node with a nil right or left.")
//		}
//		lhs, err := evaluateAst(root.left)
//		if err != nil {
//			return 0, err
//		}
//		rhs, err := evaluateAst(root.right)
//		if err != nil {
//			return 0, err
//		}
//
//		switch root.value {
//		case "+":
//			return lhs + rhs, nil
//		case "-":
//			return lhs - rhs, nil
//		case "*":
//			return lhs * rhs, nil
//		case "/":
//			return lhs / rhs, nil
//		default:
//			return 0, fmt.Errorf("Invalid operator value found for token. Value was %s but should be +, -, *, or /.", root.value)
//		}
//	}
//	return root.evaluate()
//}
//
//func Parse(input io.Reader) (int, error) {
//	nodes := make([]node, 0, 10) // default to 10 because will fit most. Could allow user provided block of mem?
//
//	// scan in nodes
//	scanner := newScanner(input)
//	for {
//		t, v := scanner.readToken()
//		if t == invalid {
//			return 0, fmt.Errorf("Got invalid token with reason: %s\n", v)
//		}
//		nodes = append(nodes, node{ t, v, nil, nil })
//
//		if t == eof {
//			break
//		}
//	}
//	ast, err := buildAstFromNodes(nodes, 0.0)
//	if err != nil {
//		return 0, err
//	}
//
//	return evaluateAst(ast)
//}
//
//
//func main() {
//  tests := []string {"(2d4 +2)\n*3"}//, "((d6+2)*3d10)/2", "(   \n(\t   10d2 +3 )\v   * 2 \n\n\n\n     )", "1dd6", "d100+b", "10d", "10dD", "2+bd100", "2+10d2b0"}
//	for _, test := range tests {
//		res, err := Parse(strings.NewReader(test))
//		if err != nil {
//			fmt.Printf("Error for input '%s'. Error was %s\n", test, err.Error())
//		} else {
//			fmt.Printf("Rolled %d for %s\n", res, test)
//		}
//	}
//}

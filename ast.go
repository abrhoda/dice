package dice

import "fmt"

type node struct {
	token token
	left  *node
	right *node
}

func walk(root *node) (int, error) {
	if root.token.kind == eof {
		return 0, nil
	}

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

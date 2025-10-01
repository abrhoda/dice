package dice

import "errors"
import "strings"
import "fmt"
import "strconv"
import "math/rand/v2"

type weight struct {
  left float64
  right float64
}

var operatorWeights = map[string]weight {
  "+": { 1.0, 1.1},
  "-": { 1.0, 1.1},
  "*": { 2.0, 2.1},
  "/": { 2.0, 2.1},
}

type tokenType int
const (
  dice tokenType = iota
  literal
  operator
  eof
)

type token struct {
   tokenType tokenType
   value string
}

func (t token) evaluate() (int, error) {
  switch t.tokenType {
    case dice:
        idx := strings.Index(t.value, "d")
        if idx == -1 {
          idx = strings.Index(t.value, "D")
          if idx == -1 {
            return 0, fmt.Errorf("Did not find 'd' or 'D' in token with type = dice and value = %s", t.value)
          }
        }

        count, err := strconv.Atoi(t.value[:idx])
        if err != nil {
          return 0, err
        }

        faces, err := strconv.Atoi(t.value[idx+1:])
        if err != nil {
          return 0, err
        }
        
        total := 0
        for range count {
          total += (rand.IntN(faces) + 1)
        }
        return total, nil

    case literal:
      return strconv.Atoi(t.value)
    default:
      return 0, fmt.Errorf("token type %d does not support evaluate.", t.tokenType)
  }
}


// functions
func isWhiteSpace(c byte) bool {
  return c == ' ' || c == '\n' || c == '\r' || c == '\v' || c == '\t'
}

func isDigit(c byte) bool {
  return c >= '0' && c <= '9'
}

func tokenize(input string) ([]token, error) {
  tokens := make([]token, 0, 10) // picking a default capacity that will fit most simple expressions
  start, end := 0, 0
  isExpression := false

  for end < len(input) {
    c := input[end]
    if isDigit(c) {
      end++
    } else if c == 'd' || c == 'D' {
      if isExpression {
        return nil, fmt.Errorf("Multiple d/D in the same expression at position %d.", end)
      }
      isExpression = true
      end++
    } else if c == '+' || c == '-' || c == '*' || c == '/' || c == '(' || c == ')' {
      if (start != end) {
        if !isExpression {
          tokens = append(tokens, token{ literal, string(input[start:end]) })
        } else {
          if (input[end-1] == 'd' || input[end-1] == 'D') {
            return nil, fmt.Errorf("Dice expression was malformed for token %s at position %d", input[start:end], end-1)
          }
          tokens = append(tokens, token{ dice, string(input[start:end]) })
          isExpression = false
        }
      }
      tokens = append(tokens, token{ operator, string(c) })
      end++
      start = end
    } else if isWhiteSpace(c) {
      if (start != end) {
        if !isExpression {
          tokens = append(tokens, token{ literal, string(input[start:end]) })
        } else {
          if (input[end-1] == 'd' || input[end-1] == 'D') {
            return nil, fmt.Errorf("Dice expression was malformed for token %s at position %d", input[start:end], end-1)
          }
          tokens = append(tokens, token{ dice, string(input[start:end]) })
          isExpression = false
        }
      }
      end++
      start = end
    } else {
      return nil, fmt.Errorf("Unknown/invalid character (%c) found at position %d", c, end)
    }
  }
  if (start != end) {
    if !isExpression {
      tokens = append(tokens, token{ literal, string(input[start:end]) })
    } else {
      if (input[end-1] == 'd' || input[end-1] == 'D') {
        return nil, fmt.Errorf("Dice expression was malformed for token %s at position %d", input[start:end], end-1)
      }
      tokens = append(tokens, token{ dice, string(input[start:end]) })
      isExpression = false
    }
  }
  // EOF token's value should be ignored.
  return append(tokens, token{ eof, "" }), nil
}

type node struct {
	token token
	left *node
	right *node
}

func nextToken(tokens []token) (token, []token) {
	return tokens[0], tokens[1:]
}

func buildAstFromTokens(tokens []token, mbp float64) (*node, error) {
	var err error
	head, tokens := nextToken(tokens)
	root := &node{ head, nil, nil }
	if root.token.tokenType == operator && root.token.value == "(" {
		root, err = buildAstFromTokens(tokens, 0.0)
		if err != nil {
			return nil, err
		}
		var temp token
		temp, tokens = nextToken(tokens)
		if temp.value != ")" {
			return nil, fmt.Errorf("Expression should have closing paren but none were found.")
		}

	} else if root.token.tokenType != dice && root.token.tokenType != literal {
		return nil, fmt.Errorf("Expression must start with a dice or literal. Found %d.", root.token.tokenType)
	}

	for {
		if len(tokens) == 0 {
			return root, nil
		}

		eofOrOp := tokens[0]
		if eofOrOp.tokenType == eof || eofOrOp.value == ")" {
			return root, nil
		} else if eofOrOp.tokenType == dice || eofOrOp.tokenType == literal {
			return nil, fmt.Errorf("Expected EOF or operation token. Found %d with value %s", eofOrOp.tokenType, eofOrOp.value)
		}
		
		lbp, rbp := operatorWeights[eofOrOp.value].left, operatorWeights[eofOrOp.value].right
		if lbp < mbp {
			break
		}

		eofOrOp, tokens = nextToken(tokens)
		rhs, err := buildAstFromTokens(tokens, rbp)

		if err != nil {
			return nil, err
		}

		root = &node{ eofOrOp, root, rhs }

	}
	
	return root, nil	
}

func evaluateAst(root *node) (int, error) {
	if root.token.tokenType == eof {
		return 0, nil
	}

	if root.token.tokenType == operator {
		if root.left == nil || root.right == nil {
			return 0, fmt.Errorf("root node is an operator node with a nil right or left.")
		}
		lhs, err := evaluateAst(root.left)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateAst(root.right)
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
// public entrypoint function for the module
func Parse(input string) (int, error) {
  if strings.TrimSpace(input) == "" {
    return 0, errors.New("Input string was empty.")
  }

  tokens, err := tokenize(input)
  if err != nil {
    return 0, err
  }

	ast, err := buildAstFromTokens(tokens, 0.0)
	if err != nil {
		return 0, err
	}

  return evaluateAst(ast)
}

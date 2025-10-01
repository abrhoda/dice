package dice

import "errors"
import "strings"
import "fmt"

// types
type weight struct {
  left float64
  right float64
}

var operatorWeights = map[byte]weight {
  '+': { 1.0, 1.1},
  '-': { 1.0, 1.1},
  '*': { 2.0, 2.1},
  '/': { 2.0, 2.1},
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

// public entrypoint function for the module
func Parse(input string) (int, error) {
  if strings.TrimSpace(input) == "" {
    return 0, errors.New("Input string was empty.")
  }

  tokens, err := tokenize(input)
  if err != nil {
    return 0, err
  }

  for idx, token := range tokens {
    fmt.Printf("%d. %s", idx, token.value)
  }
  
  return 0, nil
}

package dice

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

type tokenType int

const (
	eof tokenType = iota
	operator
	dice
	literal
)
const EOF = byte(0)

type token struct {
	kind  tokenType
	value string
}

func (token token) evaluate() (int, error) {
	switch token.kind {
	case dice:
		idx := strings.Index(token.value, "d")
		if idx == -1 {
			// NOTE there should never be a 'D' in the slice as scanner.readToken is turning this to 'd'
			idx = strings.Index(token.value, "D")
			if idx == -1 {
				return 0, fmt.Errorf("Did not find 'd' or 'D' in token with type = dice and value = %s", token.value)
			}
		}

		count := 1
		var err error

		if idx != 0 {
			count, err = strconv.Atoi(token.value[:idx])
			if err != nil {
				return 0, err
			}
		}

		faces, err := strconv.Atoi(token.value[idx+1:])
		if err != nil {
			return 0, err
		}

		total := 0
		for range count {
			total += (rand.IntN(faces) + 1)
		}
		return total, nil

	case literal:
		return strconv.Atoi(string(token.value))
	default:
		return 0, fmt.Errorf("Token type %d does not support evaluate.", token.kind)
	}
}

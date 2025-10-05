package main
import (
	"fmt"
	"math/rand/v2"
	"strings"
	"strconv"
)

type tType int
const (
	eof tType = iota
	operator
	dice
	literal
)
const EOF = byte(0)

type token struct {
	tType tType
	startPos int // start position in the 
	length int
}

// TODO changing token to have startPos and length fields instead of a value field broke this
func (token token) evaluate() (int, error) {
	switch token.tType {
	case dice:
		idx := strings.Index(token.value, "d")
		if idx == -1 {
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
		return strconv.Atoi(token.value)
	default:
		return 0, fmt.Errorf("Token type %d does not support evaluate.", token.tType)
	}
}

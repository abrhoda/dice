package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// TODO - consider not using bufio.Reader
// a lot of lexers/scanners/tokenizers use a struct like below:
//
// type scanner2 struct {
// 	buffer []byte
// 	startPos int // start position of the current token
// 	currentPos int // current position over the entire buffer
// 	endPos int // last position of the current token
// }
//
// and then token would be: 
// 
// type token2 struct {
// 	tokenType tokenTyp // eof, operator, literal, or dice 
// 	startPos // starting position of this token in the scanner's buffer
// 	length // how many bytes make up this token
// }

type scanner struct {
	buffer []byte // input string of bytes
	startPos int // start position of the current token
	currentPos int // current position over the entire buffer
 	endPos int // last position of the current token
}

type scannerError struct {
	Message string
	index int
}

func (err scannerError) Error() string {
	return err.Message
}

func newScanner(buffer []byte) *scanner {
	return &scanner{buffer, 0, 0, 0}
}

// peekByte returns the byte at currentPos without advancing the cursor
func (scanner *scanner) peekByte() byte {
	if (scanner.currentPos) >= len(scanner.buffer) {
		return EOF
	}
	return scanner.buffer[scanner.currentPos]
}

// readByte returns the byte at currentPos and advances the cursor
func (scanner *scanner) readByte() byte {
	if (scanner.currentPos) >= len(scanner.buffer) {
		return EOF
	}
	p := scanner.currentPos
	scanner.currentPos++
	return scanner.buffer[p]
}

// functions
func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\v' || b == '\t'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isDiceCharacter(b byte) bool {
	return b == 'd' || b == 'D'
}

func isOperator(b byte) bool {
  return b == '+' || b == '-' || b == '*' || b == '/' || b == '(' || b == ')'
}

// limit the bytes to a subset
func isValidByte(b byte) bool {
	return isWhiteSpace(b) || isDigit(b) || isDiceCharacter(b) || isOperator(b) || b == EOF
}

// reads the next token from the scanner's bufio.reader
// doesn't return error but instead returns (invalid, reason) on error
func (scanner *scanner) readToken() (token, error) {
	b := scanner.readByte()

	// remove this and all subsequent whitespace characters
	for isWhiteSpace(b) {
		b = scanner.readByte()
	}
	
	if !isValidByte(b) {
		return token{}, fmt.Errorf("Invalid byte (%c) found in token.", b)
	}

	if b == EOF {
		return token { eof, "" }, nil
	}
	
	if isOperator(b) {
		return &token{operator, string(b)} , nil
	}

	var sb strings.Builder
	sb.WriteByte(b)

	isDiceExp := false
	if isDiceCharacter(b) {
		isDiceExp = true
	}
	
	// peek bytes 1 by 1, checking type and, when appropriate, adding to string builder by actually reading and then peeking the next byte
	p  := scanner.peekByte()
	for {
		if isDigit(p) {
			sb.WriteByte(scanner.readByte())
			p = scanner.peekByte()
		} else if isDiceCharacter(p) {
			if isDiceExp {
				return nil, fmt.Errorf("Multiple d/D characters found in dice expression after %s.", sb.String())
			}
			isDiceExp = true

			b = scanner.readByte()
			if b == 'D' {
				b = 'd'
			}
			sb.WriteByte(b)

			// check if byte after d/D is a digit
			p = scanner.peekByte()
			if !isDigit(p) {
				return nil, fmt.Errorf("Character after d/D not a digit. Found %c", p)
			}
		} else if isWhiteSpace(p) || isOperator(p) || p == EOF {
			break
		} else {
			return nil, fmt.Errorf("Invalid byte (%c) found in token.", p)
		}
	}

	if isDiceExp {
		return &token{ dice, sb.String() }, nil
	} else {
		return &token{ literal, sb.String() }, nil
	}
}

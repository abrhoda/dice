package main

import (
	"fmt"
)

type scanner struct {
	buffer     []byte // input string of bytes
	startPos   int    // start position of the current token
	currentPos int    // current position over the entire buffer
	//endPos int // last position of the current token
}

type scannerError struct {
	Message string
	index   int
}

func (err scannerError) Error() string {
	return err.Message
}

func newScanner(buffer []byte) *scanner {
	return &scanner{buffer, 0, 0}
}

// peekByte returns the byte at currentPos without advancing the cursor
func (scanner *scanner) peekByte() byte {
	if (scanner.currentPos) < len(scanner.buffer) {
		return scanner.buffer[scanner.currentPos]
	}
	return EOF
}

// readByte returns the byte at currentPos and advances the cursor
func (scanner *scanner) readByte() byte {
	if (scanner.currentPos) < len(scanner.buffer) {
		p := scanner.currentPos
		scanner.currentPos++
		return scanner.buffer[p]
	}
	return EOF
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
		scanner.startPos++
		b = scanner.readByte()
	}

	if !isValidByte(b) {
		return token{}, scannerError{"Invalid byte (%c) found in buffer.", scanner.currentPos}
	}

	if b == EOF {
		t := token{eof, string(scanner.buffer[scanner.startPos:scanner.currentPos])}
		scanner.startPos = scanner.currentPos
		return t, nil
	}

	if isOperator(b) {
		t := token{operator, string(scanner.buffer[scanner.startPos:scanner.currentPos])}
		scanner.startPos = scanner.currentPos
		return t, nil
	}

	isDiceExp := false
	if isDiceCharacter(b) {
		isDiceExp = true
	}

	// peek bytes 1 by 1, checking type and, when appropriate, adding to string builder by actually reading and then peeking the next byte
	p := scanner.peekByte()
	for {
		if isDigit(p) {
			_ = scanner.readByte()
			p = scanner.peekByte()
		} else if isDiceCharacter(p) {

			// NOTE this isn't needed.
			//if isDiceExp {
			//	return token{}, scannerError{"Multiple d/D characters found in dice expression.", scanner.currentPos}
			//}

			isDiceExp = true

			b = scanner.readByte()
			if b == 'D' {
				b = 'd'
			}
			// check if byte after d/D is a digit
			p = scanner.peekByte()
			if !isDigit(p) {
				return token{}, scannerError{fmt.Sprintf("Character after d/D not a digit. Found %c", p), scanner.currentPos}
			}
		} else if isWhiteSpace(p) || isOperator(p) || p == EOF {
			break
		} else {
			return token{}, scannerError{fmt.Sprintf("Invalid byte (%c) found in token.", p), scanner.currentPos}
		}
	}

	// TODO clean this up. Default and then check and change isn't great.
	t := token{literal, string(scanner.buffer[scanner.startPos:scanner.currentPos])}
	if isDiceExp {
		t.kind = dice
	}
	scanner.startPos = scanner.currentPos
	return t, nil
}

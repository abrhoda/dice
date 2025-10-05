package main

import "fmt"

func main() {
	tests := []string{"(2d4 +2)\n*3", "((d6+2)*3d10)/2", "(   \n(\t   10d2 +3 )\v   * 2 \n\n\n\n     )", "1dd6", "d100+b", "10d", "10dD", "2+bd100", "2+10d2b0"}
	for _, test := range tests {
		buffer := []byte(test)
		s := newScanner(buffer)
		for {

			t, err := s.readToken()
			if err != nil {
				fmt.Printf("Error for input '%s'. Error was %s\n", test, err.Error())
			}

			if t.tType == eof {
				break
			}

			fmt.Printf("Token type = %d, starting at %d and len of %d. value was '%s'\n", t.tType, t.startPos, (t.endPos - t.startPos), string(buffer[t.startPos:t.endPos]))
		}
	}
}

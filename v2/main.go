package main

import "fmt"

func main() {
	tests := []string{"(2d4 +2)\n*3", "((d6+2)*3d10)/2", "(   \n(\t   10d2 +3 )\v   * 2 \n\n\n\n     )", "1dd6", "d100+b", "10d", "10dD", "2+bd100", "2+10d2b0"}
	//tests := []string{"1dd6", "d100+b", "10d", "10dD", "2+bd100", "2+10d2b0"}
	for _, test := range tests {
		fmt.Println(test)
		buffer := []byte(test)
		p := NewParser(buffer)
		err := p.Parse()
		if err != nil {
			fmt.Printf("Error happened: %s\n", err.Error())
			continue
		}

		ast, err := p.astFromTokens(0.0)
		if err != nil {
			fmt.Printf("Error happened: %s\n", err.Error())
			continue
		}

		res, err := walk(ast)
		if err != nil {
			fmt.Printf("Error happened: %s\n", err.Error())
			continue
		}

		fmt.Printf("return = %d\n", res)
	}
}

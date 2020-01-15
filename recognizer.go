package main

import (
	"fmt"
	"lexer"
	"os"
	"parser"
)

func main() {
	var lxr *lexer.Lexer
	if len(os.Args) > 1 {
		lxr = lexer.NewFromFileName(os.Args[1])
	} else {
		lxr = lexer.NewFromFile(os.Stdin)
	}

	psr := parser.New(lxr)

	r := psr.Recognizer()

	if r {
		fmt.Printf("It's an expression\n")
	} else {
		fmt.Printf("It's NOT an expression\n")
	}
}

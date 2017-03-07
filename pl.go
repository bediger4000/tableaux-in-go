package main

import (
	"os"
	"fmt"
	"lexer"
)

func main() {
	var lxr *lexer.Lexer
	if len(os.Args) > 1 {
		lxr = lexer.NewFromFileName(os.Args[1])
	} else {
		lxr = lexer.NewFromFile(os.Stdin)
	}

	tokenType, token := lxr.NextToken()

	fmt.Printf("Token %q, type %d\n", token, tokenType)

}

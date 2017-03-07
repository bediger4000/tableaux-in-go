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

	for token, tokenType := lxr.NextToken(); tokenType != lexer.EOF; token, tokenType = lxr.NextToken() {
		fmt.Printf("Token %q, type %s, %d\n", token, lexer.TokenName(tokenType), tokenType)
	}

}

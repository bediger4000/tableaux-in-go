package main

import (
	"fmt"
	"lexer"
	"node"
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

	var root *node.Node
	root = psr.Parse()

	if root != nil {
		root.Print(os.Stdout)
		fmt.Printf("\n")
	}
}

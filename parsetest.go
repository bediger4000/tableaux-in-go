package main

import (
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

	var tree *node.Node
	tree = psr.Parse()

	tree.Print(os.Stdout)
}

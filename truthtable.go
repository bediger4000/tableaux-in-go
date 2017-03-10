package main

import (
	"fmt"
	"lexer"
	"node"
	"os"
	"parser"
	"sort"
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

	printTruthTable(root)

	os.Exit(0)
}

func printTruthTable(root *node.Node) {

	identifiers := findIdentifiers(root)

	fmt.Printf("Found %d identifiers: %v\n", len(identifiers), identifiers)
}

func findIdentifiers(n *node.Node) ([]string) {

	allIdentifiers := findAllIdentifiers(n)

	alreadySeen := make(map[string]bool)
	var uniqIdentifiers []string
	for _, id := range allIdentifiers {
		if _, ok := alreadySeen[id]; !ok {
			alreadySeen[id] = true
			uniqIdentifiers = append(uniqIdentifiers, id)
		}
	}

	sort.Strings(uniqIdentifiers)

	return uniqIdentifiers
}

func findAllIdentifiers(n *node.Node) ([]string) {

	var identifiers []string
	if n.Op == lexer.IDENT {
		identifiers = append(identifiers, n.Ident)
	}
	if n.Left != nil {
		ids := findAllIdentifiers(n.Left)
		identifiers = append(identifiers, ids...)
	}
	if n.Right != nil {
		ids := findAllIdentifiers(n.Right)
		identifiers = append(identifiers, ids...)
	}
	return identifiers
}

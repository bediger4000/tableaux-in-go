package main

import (
	"fmt"
	"lexer"
	"log"
	"node"
	"os"
	"parser"
	"sort"
	"stringbuffer"
	"tableaux"
)

func main() {
	var lxr *lexer.Lexer
	if len(os.Args) > 1 {
		var expr stringbuffer.Buffer
		expr.Store(os.Args[1] + "\n")
		lxr = lexer.NewFromFile(&expr)
	} else {
		lxr = lexer.NewFromFile(os.Stdin)
	}

	psr := parser.New(lxr)

	var root *node.Node
	root = psr.Parse()
	fmt.Printf("Expression: %q\n", node.ExpressionToString(root))

	spork := tableaux.New(root, false, nil)

	spork.DoTnode()
	spork.PrintTnode()

	os.Exit(0)
}

func printTruthTable(root *node.Node) {

	identifiers := findIdentifiers(root)

	max := len(identifiers)
	var vals []bool
	for i := 0; i < max; i++ {
		vals = append(vals, true)
	}

	for _, variable := range identifiers {
		n := 5 - len(variable)
		spacer := ""
		for i := 0; i < n; i++ {
			spacer += " "
		}
		fmt.Printf("%s%s ", spacer, variable)
	}
	expression := node.ExpressionToString(root)
	fmt.Printf("\t%s\n", expression)

	for {

		valuation := make(map[string]bool)
		for idx, id := range identifiers {
			valuation[id] = vals[idx]
		}
		r := evaluateExpression(root, valuation)
		printRow(identifiers, vals, r)

		var idx int
		for idx = max - 1; idx >= 0; idx-- {
			vals[idx] = !vals[idx]
			if !vals[idx] {
				break
			}
		}
		if idx < 0 {
			break
		}
	}
}

func findIdentifiers(n *node.Node) []string {

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

func findAllIdentifiers(n *node.Node) []string {

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

func evaluateExpression(n *node.Node, valuation map[string]bool) bool {
	switch n.Op {
	case lexer.NOT:
		return !evaluateExpression(n.Left, valuation)
	case lexer.AND:
		return evaluateExpression(n.Left, valuation) && evaluateExpression(n.Right, valuation)
	case lexer.OR:
		return evaluateExpression(n.Left, valuation) || evaluateExpression(n.Right, valuation)
	case lexer.IMPLIES:
		p := evaluateExpression(n.Left, valuation)
		q := evaluateExpression(n.Right, valuation)
		if !p && q {
			return false
		}
		return true
	case lexer.EQUIV:
		return evaluateExpression(n.Left, valuation) == evaluateExpression(n.Right, valuation)
	case lexer.IDENT:
		return valuation[n.Ident]
	}
	log.Fatalf("Problem with node type %s (%d): shouldn't get here\n", lexer.TokenName(n.Op), n.Op)
	return false
}

func printRow(identifiers []string, vals []bool, r bool) {
	for idx, _ := range identifiers {
		spacer := " "
		if !vals[idx] {
			spacer = ""
		}
		fmt.Printf("%s%v ", spacer, vals[idx])
	}
	fmt.Printf("\t%v\n", r)
}

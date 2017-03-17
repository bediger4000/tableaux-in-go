package main

import (
	"bytes"
	"flag"
	"fmt"
	"lexer"
	"log"
	"node"
	"os"
	"parser"
	"sort"
	"tableaux"
)

func main() {

	graphVizOutputFilename := flag.String("g", "", "File name for graphviz output, no default")
	flag.Parse()

	var expressions []string

	if flag.NArg() > 0 {
		expressions = flag.Args()
	}

	if len(expressions) == 0 {
		fmt.Fprintf(os.Stderr, "Need at least one propostitional logic formula on command line\n")
		os.Exit(1)
	}

	// Parse expression(s) on cmd line into *node.Node objects
	var trees []*node.Node

	for _, expression := range expressions {
		var lxr *lexer.Lexer
		expr := bytes.NewBufferString(expression+"\n")  // parser.Parser needs to recognize end-of-line
		lxr = lexer.NewFromFile(expr)
		psr := parser.New(lxr)
		tree := psr.Parse()
		fmt.Printf("Expression: %q\n", node.ExpressionToString(tree))
		trees = append(trees, tree)
	}

	// tblx will become the entire tableau, below
	tblx := tableaux.New(trees[0], false, nil)

	if len(trees) > 1 {
		// More than 1 PL formula, put them together for deciding
		// logical consequence - all signed T except that last one F.
		var t *tableaux.Tnode
		for _, tree := range trees {
			t = tableaux.New(tree, true, nil)  // All signed T
			tblx.AppendLeaf(t)
		}
		t.Sign = false  // Except final one, signed F
	} else {
		// Single expression. Subjoin its own inferences.
		tblx.AddInferences(tblx)
		tblx.Used = true
	}
	
	tautological := false  // The answer we're looking for.
	foundUnused  := true   // Found a formula with no previously subjoined inferences

	for foundUnused {
		unclosedLeaves := tblx.FindUnclosedLeaf()
		if len(unclosedLeaves) == 0 {
			fmt.Printf("No unclosed branches\n")
			tautological = true
			break
		}
		foundUnused = false
		for _, leaf := range unclosedLeaves {

fmt.Printf("Unclosed leaf, %v: %q\n", leaf.Sign, leaf.Expression)
			unusedFormula := leaf.FindTallestUnused()
			if unusedFormula != nil {
fmt.Printf("Unused formula above leaf: %v: %q\n", unusedFormula.Sign, unusedFormula.Expression)
				uncl := unusedFormula.FindUnclosedLeaf()
				// uncl will have leaf as an element, but subjoin inferences
				// to all leaf nodes under unusedFormula
				for _, leafNode := range uncl {
					leafNode.AddInferences(unusedFormula)
				}
				foundUnused = true
				unusedFormula.Used = true
				// Just subjoined inferences to each unclosed leaf node
				// in the branch below unusedFormula, so uncl and unclosedLeaves
				// now don't have all unclosed leaf nodes in them.
				break
			}
		}
	}

	if tautological {
		fmt.Printf("Formula is tautology\n")
	} else {
		fmt.Printf("Formula is not a tautology\n")
	}

	fmt.Printf("\n*/\n")

	if *graphVizOutputFilename != "" {
		fout, err := os.OpenFile(*graphVizOutputFilename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("Problem opening %d write-only: %s\n", *graphVizOutputFilename, err)
			os.Exit(1)
		}
		defer fout.Close()
		tblx.GraphTnode(fout)
	}

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

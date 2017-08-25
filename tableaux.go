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
		fmt.Fprintf(os.Stderr, "Need at least one propositional logic formula on command line\n")
		os.Exit(1)
	}

	// Parse expression(s) on cmd line into *node.Node objects
	var trees []*node.Node

	expressionCount := len(expressions)
	denotation := "Expression"
	if expressionCount > 1 {
		denotation = "Hypothesis"
	}

	for idx, expression := range expressions {
		var lxr *lexer.Lexer
		expr := bytes.NewBufferString(expression + "\n") // parser.Parser needs to recognize end-of-line
		lxr = lexer.NewFromFile(expr)
		psr := parser.New(lxr)
		tree := psr.Parse()
		fmt.Printf("%s: %q\n", denotation, node.ExpressionToString(tree))
		trees = append(trees, tree)
		if idx == expressionCount-2 {
			denotation = "Consequence"
		}
	}

	// tblx will become the entire tableau, below
	tblx := tableaux.New(trees[0], false, nil)

	var finalFormula *tableaux.Tnode
	if len(trees) > 1 {
		// More than 1 PL formula, put them together for deciding
		// logical consequence - all signed T except that last one F.
		tblx.Sign = true
		for _, tree := range trees[1:] {
			finalFormula = tableaux.New(tree, true, nil) // All signed T
			tblx.AppendLeaf(finalFormula)
		}
		finalFormula.Sign = false // Except final one, signed F
	} else {
		// Single expression. Subjoin its own inferences.
		tblx.AddInferences(tblx)
		tblx.Used = true
	}

	tautological := false // The answer we're looking for.
	foundUnused := true   // Found a formula with no previously subjoined inferences

	for foundUnused {
		unclosedLeaves := tblx.FindUnclosedLeaf()
		if len(unclosedLeaves) == 0 {
			tautological = true
			break
		}
		foundUnused = false
		for _, leaf := range unclosedLeaves {

			unusedFormula := leaf.FindTallestUnused()
			if unusedFormula != nil {
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

	fmt.Printf("/*\n")

	tableaux.PrintTableaux(os.Stdout, tblx)

	var modifier string
	if !tautological {
		modifier = " not"
	} else {
		modifier = ""
	}

	if finalFormula == nil {
		fmt.Printf("Formula is%s a tautology\n", modifier)
	} else {
		fmt.Printf("%s is%s a logical consequence of hypotheses\n", finalFormula.Expression, modifier)
	}

	fmt.Printf("*/\n")

	if *graphVizOutputFilename != "" {
		fout, err := os.OpenFile(*graphVizOutputFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("Problem opening %q write-only: %s\n", *graphVizOutputFilename, err)
			os.Exit(1)
		}
		defer fout.Close()
		tblx.GraphTnode(fout)
	}
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"lexer"
	"node"
	"os"
	"parser"
)

func main() {

	graphVizOutputFilename := flag.String("g", "", "File name for graphviz output, no default")
	flag.Parse()

	var lxr *lexer.Lexer
	if flag.NArg() > 0 {
		expressions := flag.Args()
		expr := bytes.NewBufferString(expressions[0] + "\n")
		lxr = lexer.NewFromFile(expr)
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

	if *graphVizOutputFilename != "" {
		fout, err := os.OpenFile(*graphVizOutputFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem opening %q write-only: %s\n", *graphVizOutputFilename, err)
			os.Exit(1)
		}
		defer fout.Close()
		root.GraphNode(fout)
	}
}

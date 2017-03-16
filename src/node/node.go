package node

// Parse tree - a binary tree of objects of type Node,
// and associated utility functions and methods.

import (
	"bytes"
	"fmt"
	"io"
	"lexer"
)

// All elements exported, everything reaches inside instances of Node
// to find things out, or to change Left and Right. Private elements
// would cost me gross ol' getter and setter boilerplate.
type Node struct {
	Op    lexer.TokenType
	Ident string
	Left  *Node
	Right *Node
}

// Create interior nodes of a parse tree, which will
// all have a &, ~, |, >, = operator associated.
func NewOpNode(op lexer.TokenType) *Node {
	var n Node
	n.Op = op
	return &n
}

// Create leaf nodes of a parse tree, which should all
// be lexer.IDENT identifier nodes.
func NewIdentNode(identifier string) *Node {
	var n Node
	n.Op = lexer.IDENT
	n.Ident = identifier
	return &n
}

// Put a human-readable, nicely formatted string representation
// of a parse tree onto the io.Writer, w.
// Essentially just an in-order traversal of a binary tree, with
// accomodating a few oddities, like parenthesization, and the
// "~" (not) operator being a prefix.
func (p *Node) Print(w io.Writer) {

	if lexer.NOT == p.Op {
		fmt.Fprintf(w, "~")
	}
	if p.Left != nil {
		printParen := false
		if p.Left.Op != lexer.IDENT && p.Left.Op != lexer.NOT {
			fmt.Fprintf(w, "(")
			printParen = true
		}
		p.Left.Print(w)
		if printParen {
			fmt.Fprintf(w, ")")
		}
	}

	var oper rune
	switch p.Op {
	case lexer.IMPLIES:
		oper = '>'
		break
	case lexer.AND:
		oper = '&'
		break
	case lexer.OR:
		oper = '|'
		break
	case lexer.EQUIV:
		oper = '='
		break
	}
	if oper != 0 {
		fmt.Fprintf(w, " %c ", oper)
	}

	if p.Op == lexer.IDENT {
		fmt.Fprintf(w, "%s", p.Ident)
	}

	if p.Right != nil {
		printParen := false
		if p.Right.Op != lexer.IDENT && p.Right.Op != lexer.NOT {
			fmt.Fprintf(w, "(")
			printParen = true
		}
		p.Right.Print(w)
		if printParen {
			fmt.Fprintf(w, ")")
		}
	}
}

// Creating a Golang string with a human readable representation
// of a parse tree in it.
func ExpressionToString(root *Node) (string) {
	var sb bytes.Buffer
	root.Print(&sb)
	return sb.String()
}

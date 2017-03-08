package node

import (
	"fmt"
	"io"
	"lexer"
)

type Node struct {
	op    lexer.TokenType
	ident string
	Left  *Node
	Right *Node
}

func NewOpNode(op lexer.TokenType) *Node {
	var n Node
	n.op = op
	return &n
}

func NewIdentNode(identifier string) *Node {
	var n Node
	n.op = lexer.IDENT
	n.ident = identifier
	return &n
}

func (p *Node) Print(w io.Writer) {

	if lexer.NOT == p.op {
		fmt.Fprintf(w, "~")
	}
	if p.Left != nil {
		printParen := false
		if p.Left.op != lexer.IDENT && p.Left.op != lexer.NOT {
			fmt.Fprintf(w, "(")
			printParen = true
		}
		p.Left.Print(w)
		if printParen {
			fmt.Fprintf(w, ")")
		}
	}

	var oper rune
	switch p.op {
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

	if p.op == lexer.IDENT {
		fmt.Fprintf(w, "%s", p.ident)
	}

	if p.Right != nil {
		printParen := false
		if p.Right.op != lexer.IDENT && p.Right.op != lexer.NOT {
			fmt.Fprintf(w, "(")
			printParen = true
		}
		p.Right.Print(w)
		if printParen {
			fmt.Fprintf(w, ")")
		}
	}
}

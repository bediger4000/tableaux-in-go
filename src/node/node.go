package node

import (
	"fmt"
	"io"
	"lexer"
	"stringbuffer"
)

type Node struct {
	Op    lexer.TokenType
	Ident string
	Left  *Node
	Right *Node
}

func NewOpNode(op lexer.TokenType) *Node {
	var n Node
	n.Op = op
	return &n
}

func NewIdentNode(identifier string) *Node {
	var n Node
	n.Op = lexer.IDENT
	n.Ident = identifier
	return &n
}

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

func ExpressionToString(root *Node) (string) {
	var sb stringbuffer.Buffer
	root.Print(&sb)
	return sb.String()
}

package parser

import (
	"fmt"
	"lexer"
	"node"
	"os"
)

// Parser instances get used to do a parse of a single
// logical expression. Since instances of Lexer have a
// io.Reader or a string in them, Parser instances have
// no idea what they're parsing from.
type Parser struct {
	lexer *lexer.Lexer
}

var nextOp [9]lexer.TokenType

// New used to create a Parser instance, injecting
// a prepared Lexer instance.
func New(lxr *lexer.Lexer) *Parser {
	var parser Parser
	parser.lexer = lxr
	nextOp[int(lexer.EQUIV)] = lexer.IMPLIES
	nextOp[int(lexer.IMPLIES)] = lexer.OR
	nextOp[int(lexer.OR)] = lexer.AND
	nextOp[int(lexer.AND)] = lexer.EQUIV
	return &parser
}

// Parse creates a parse tree in the form of a
// binary tree of pointers to node.Node, from
// whatever source of text the Lexer instance
// has in it.
func (p *Parser) Parse() *node.Node {
	root := p.parseProduction(lexer.EQUIV)
	if root != nil {
		q := p.expect(lexer.EOL)
		if !q {
			root = nil
		}
	}
	return root
}

// See README.md: basically 4 of the 5 productions look like:
// Nonterminal0 -> Nonterminal1 {op1 Nonterminal1}
// Nonterminal1 -> Nonterminal2 {op2 Nonterminal2}
//  ...
// The code for each parsing method was almost identical, except
// for the next function to call, and the condition on the for-loop.
// Generalize all 4 of the parseNonterminal() methods into one method.

func (p *Parser) parseProduction(op lexer.TokenType) *node.Node {

	nextProduction := p.parseProduction
	if op == lexer.AND {
		nextProduction = p.parseFactor
	}

	no := nextOp[op]
	newNode := nextProduction(no) // Weird that this works.
	if newNode != nil {
		for _, typ := p.lexer.Next(); typ == op; _, typ = p.lexer.Next() {
			p.lexer.Consume()
			tmp := node.NewOpNode(op)
			tmp.Left = newNode
			tmp.Right = nextProduction(no) // p.parseProduction(no) or p.parseFactor(no)
			newNode = tmp
		}
	}
	return newNode
}

func (p *Parser) parseFactor(op lexer.TokenType) *node.Node {
	var n *node.Node

	token, typ := p.lexer.Next()

	switch typ {
	case lexer.IDENT:
		p.lexer.Consume()
		n = node.NewIdentNode(token)
	case lexer.LPAREN:
		p.lexer.Consume()
		n = p.parseProduction(op)
		if n != nil {
			if !p.expect(lexer.RPAREN) {
				fmt.Fprintf(os.Stderr, "Didn't find a right paren to match left parenthese\n")
				n = nil
			}
		}
	case lexer.NOT:
		p.lexer.Consume()
		n = node.NewOpNode(lexer.NOT)
		n.Left = p.parseFactor(op)
	default:
		fmt.Fprintf(os.Stderr, "Found token %q, type %s (%d) instead of IDENT|LPAREN|NOT\n", token, lexer.TokenName(typ), typ)
		n = nil
	}
	return n
}

func (p *Parser) expect(expectedType lexer.TokenType) bool {
	token, tokenType := p.lexer.Next()
	if tokenType == expectedType {
		p.lexer.Consume()
	} else {
		fmt.Fprintf(os.Stderr, "Expected token type %s, found %s (%q)\n", lexer.TokenName(expectedType), lexer.TokenName(tokenType), token)
		return false
	}
	return true
}

package parser

import (
	"fmt"
	"lexer"
	"os"
	"node"
)

type Parser struct {
	lexer *lexer.Lexer
}

func New(lexer *lexer.Lexer) *Parser {
	var parser Parser
	parser.lexer = lexer
	return &parser
}

func (p *Parser) Parse() (*node.Node) {
	root := p.parseEquivalence()
	if root != nil {
		q := p.Expect(lexer.EOL)
		if !q {
			root = nil
		}
	}
	return root
}

func (p *Parser) parseEquivalence() (*node.Node) {
	n := p.parseImplication()
	if n != nil {
		for _, typ := p.lexer.Next(); typ == lexer.EQUIV; _, typ = p.lexer.Next() {
			p.lexer.Consume()
			tmp := node.NewOpNode(lexer.EQUIV)
			tmp.Left = n
			tmp.Right = p.parseImplication()
			n = tmp
		}
	}
	return n
}

func (p *Parser) parseImplication() (*node.Node) {
	n := p.parseDisjunction()
	if n != nil {
		for _, typ := p.lexer.Next(); typ == lexer.IMPLIES; _, typ = p.lexer.Next() {
			p.lexer.Consume()
			tmp := node.NewOpNode(lexer.IMPLIES)
			tmp.Left = n
			tmp.Right = p.parseDisjunction()
			n = tmp
		}
	}
	return n
}

func (p *Parser) parseDisjunction() (*node.Node) {
	n := p.parseConjunction()
	if n != nil {
		for _, typ := p.lexer.Next(); typ == lexer.OR; _, typ = p.lexer.Next() {
			p.lexer.Consume()
			tmp := node.NewOpNode(lexer.OR)
			tmp.Left = n
			tmp.Right = p.parseConjunction()
			n = tmp
		}
	}
	return n
}

func (p *Parser) parseConjunction() (*node.Node) {
	n := p.parseFactor()
	if n != nil {
		for _, typ := p.lexer.Next(); typ == lexer.AND; _, typ = p.lexer.Next() {
			p.lexer.Consume()
			tmp := node.NewOpNode(lexer.AND)
			tmp.Left = n
			tmp.Right = p.parseFactor()
			n = tmp
		}
	}
	return n
}

func (p *Parser) parseFactor() (*node.Node) {
	var n *node.Node

	token, typ := p.lexer.Next()

	switch typ {
	case lexer.IDENT:
		p.lexer.Consume()
		n = node.NewIdentNode(token)
	case lexer.LPAREN:
		p.lexer.Consume()
		n = p.parseEquivalence()
		if n != nil {
			if !p.Expect(lexer.RPAREN) {
				fmt.Fprintf(os.Stderr, "Didn't find a right paren to match left parenthese\n")
				n = nil
			}
		}
	case lexer.NOT:
		p.lexer.Consume()
		n = node.NewOpNode(lexer.NOT)
		n.Left = p.parseFactor()
	default:
		fmt.Fprintf(os.Stderr, "Found token %q, type %s (%d) instead of IDENT|LPAREN|NOT\n", token, lexer.TokenName(typ), typ)
		n = nil
	}
	return n
}

func (p *Parser) Expect(expectedType lexer.TokenType) (bool) {
	token, tokenType := p.lexer.Next()
	if tokenType == expectedType {
		p.lexer.Consume()
	} else {
		fmt.Fprintf(os.Stderr, "Expected token type %s, found %s (%q)\n", lexer.TokenName(expectedType),  lexer.TokenName(tokenType),  token)
		return false
	}
	return true
}

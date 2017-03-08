package parser

import (
	"fmt"
	"lexer"
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

func (p *Parser) Parse() (n *node.Node) {

	token, typ := p.lexer.NextToken()

	fmt.Printf("%q type %s\n", token, lexer.TokenName(typ))

	switch typ {
	case lexer.NOT:
		n = node.NewOpNode(lexer.NOT)
		n.Left = p.Parse()
	case lexer.LPAREN:
		left := p.Parse()
		n = p.Parse()
		n.Left = left
		n.Right = p.Parse()
		_, _ = p.lexer.NextToken() // should be LPAREN
	case lexer.AND, lexer.OR, lexer.IMPLIES, lexer.EQUIV:
		n = node.NewOpNode(typ)
	case lexer.IDENT:
		n = node.NewIdentNode(token)
	}
	return n
}

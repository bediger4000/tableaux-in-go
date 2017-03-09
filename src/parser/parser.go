package parser

import (
	"fmt"
	"lexer"
	"log"
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
	var n node.Node
	fmt.Printf("Parse()\n")
	return &n
}

func (p *Parser) Expect(expectedType lexer.TokenType) (bool) {
	token, tokenType := p.lexer.Next()
	if tokenType == expectedType {
		p.lexer.Consume()
	} else {
		log.Printf("Expected token type %s, found %s (%q)\n", lexer.TokenName(expectedType),  lexer.TokenName(tokenType),  token)
		return false
	}
	return true
}

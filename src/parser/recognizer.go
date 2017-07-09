package parser

import (
	"lexer"
	"log"
)

// Recognizer returns true if the lexer feeds it a propositional
// logic expression, false otherwise.
func (p *Parser) Recognizer() bool {
	r := p.recognizeE()
	if r {
		q := p.expect(lexer.EOL)
		r = r && q
	}
	return r
}

func (p *Parser) recognizeE() bool {
	r := p.recognizeP()
	if r {
		for _, typ := p.lexer.Next(); lexer.BinaryOperator(typ); _, typ = p.lexer.Next() {
			p.lexer.Consume()
			r = r && p.recognizeP()
		}
	}
	return r
}

func (p *Parser) recognizeP() bool {
	r := false
	token, typ := p.lexer.Next()
	switch typ {
	case lexer.IDENT:
		p.lexer.Consume()
		r = true
	case lexer.LPAREN:
		p.lexer.Consume()
		r = p.recognizeE()
		r = p.expect(lexer.RPAREN) && r
	case lexer.NOT:
		p.lexer.Consume()
		r = p.recognizeP()
	default:
		log.Printf("Found token %q, type %s (%d) instead of IDENT|LPAREN|NOT\n", token, lexer.TokenName(typ), typ)
	}
	return r
}

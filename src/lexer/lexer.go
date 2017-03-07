package lexer

import (
	"bufio"
	"log"
	"os"
	"fmt"
	"unicode/utf8"
)

type Lexer struct {
	fileName    string
	fd      *os.File
	scanner *bufio.Scanner

	pos int
	line string
	lineLength int
}

type TokenType int

const (
	NOT  TokenType = iota
	AND  TokenType = iota
	OR   TokenType = iota
	IMPLIES   TokenType = iota
	EQUIV   TokenType = iota
	IDENT   TokenType = iota
	LPAREN   TokenType = iota
	RPAREN   TokenType = iota
	EOL   TokenType = iota
	EOF   TokenType = iota
)

func NewFromFile(file *os.File) *Lexer {
	var z Lexer
	z.fileName = "stdin"
	z.fd = file
	z.scanner = bufio.NewScanner(z.fd)
	return &z
}

func NewFromFileName(fileName string) *Lexer {
	var z Lexer
	var err error
	z.fd, err = os.Open(fileName)
	if err != nil {
		log.Fatalf("Opening file %q: %s\n", fileName, err)
	}
	z.fileName = fileName
	z.scanner = bufio.NewScanner(z.fd)
	return &z
}

func (p *Lexer) NextToken() (string, TokenType) {

	if p.lineLength > 0 && p.pos == p.lineLength {
		p.pos = -1
		p.lineLength = 0
		return "", EOL
	}

	for p.lineLength == 0 {
		// Read in next line
		if p.scanner.Scan() {
			p.line = p.scanner.Text()
			p.pos = 0
			p.lineLength = len(p.line)
		} else {
			err := p.scanner.Err()
			if err != nil {
				if err := p.scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "Reading %s: %s\n", p.fileName, err)
					return err.Error(), EOF
				}
			} else {
					return "", EOF
			}
		}
	}

	var token []rune
	var typ TokenType
	foundToken := false

	for !foundToken && p.pos < p.lineLength {
		var c rune
		c, w := utf8.DecodeRuneInString(p.line[p.pos:])

		switch c {
		case '(', ')', '&', '~', '|', '=':
			if len(token) == 0 {
				switch c {
				case '(':
					typ = LPAREN
				case ')':
					typ = LPAREN
				case '&':
					typ = AND
				case '|':
					typ = OR
				case '=':
					typ = EQUIV
				case '~':
					typ = NOT
				}
				foundToken = true
				token = append(token, c)
				p.pos += w
			}
			foundToken = true
		default:
			if c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
				typ = IDENT
				token = append(token, c)
				p.pos += w
			}
		case ' ', '\t':
			if len(token) > 0 {
				foundToken = true
			}
			p.pos += w
		}
	}

	return string(token), typ
}

func TokenName(t TokenType) (string) {
	var r string = "unknown"
	switch t {
	case LPAREN:
		r = "LPAREN"
	case RPAREN:
		r = "RPAREN"
	case NOT:
		r = "NOT"
	case AND:
		r = "AND"
	case OR:
		r = "OR"
	case IMPLIES:
		r = "IMPLIES"
	case EQUIV:
		r = "EQUIV"
	case IDENT:
		r = "IDENT"
	case EOL:
		r = "EOL"
	case EOF:
		r = "EOF"
	}
	return r
}

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
	z.scanner.Split(plSplitter)
	return &z
}

func (p *Lexer) NextToken() (string, TokenType) {

	if !p.scanner.Scan() {
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

	var typ TokenType

	token := p.scanner.Text()

	// This is kind bunk, as plSplitter() knows perfectly well what
	// type the token had, but unless I use a package-level variable,
	// I can't figure out how to communicate token type from plSplitter()
	switch token {
	case "(":
		typ = LPAREN
	case ")":
		typ = RPAREN
	case "&":
		typ = AND
	case "|":
		typ = OR
	case ">":
		typ = IMPLIES
	case "=":
		typ = EQUIV
	case "\n":
		typ = EOL
	default:
		typ = IDENT
	}

	return token, typ
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

func plSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {

	foundToken := false

	for !foundToken && advance < len(data) {
		var c rune
		c, w := utf8.DecodeRune(data[advance:])

		switch c {
		case '(', ')', '&', '~', '|', '=', '>':
			if len(token) == 0 {
				end := advance + w
				token = append(token, data[advance:end]...)
				advance = end
			}
			foundToken = true
		case ' ', '\t':
			if len(token) > 0 {
				foundToken = true
			}
			advance += w
		case '\n':
			if len(token) == 0 {
				end := advance + w
				token = append(token, data[advance:end]...)
				advance = end
			}
			foundToken = true
		default:
			if c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
				end := advance + w
				token = append(token, data[advance:end]...)
				advance = end
			} else {
				// Skip over meaningless characters
				advance += w
			}
		}
	}
	return
}

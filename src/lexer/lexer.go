package lexer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"
)

// Lexer instances give back lexemes (instances of TokenType)
// on successive calls to lexer.Next()
type Lexer struct {
	fileName     string
	fd           io.Reader
	scanner      *bufio.Scanner
	currentToken string
	currentType  TokenType
	needsRefresh bool
}

// TokenType - the "part of speech" of the propositional logic
// token/lexeme the instance of Lexer just found. Exported because
// Parser instances need to know, and these constants get used to
// denote what kind of node a piece of a parse tree is.
type TokenType int

// All the lexemes that this program knows about.
const (
	NOT     TokenType = iota
	AND     TokenType = iota
	OR      TokenType = iota
	IMPLIES TokenType = iota
	EQUIV   TokenType = iota
	IDENT   TokenType = iota
	LPAREN  TokenType = iota
	RPAREN  TokenType = iota
	EOL     TokenType = iota
	EOF     TokenType = iota
)

// NewFromFile creates a lexer that reads text from an io.Reader
// when finding lexemes. The io.Reader comes from a file, or from
// an instance of bytes.Buffer, which just holds a string.
func NewFromFile(file io.Reader) *Lexer {
	var z Lexer
	z.fileName = "stdin"
	z.fd = file
	z.scanner = bufio.NewScanner(z.fd)
	z.scanner.Split(plSplitter)
	z.needsRefresh = true
	return &z
}

// NewFromFileName conveniently gives back a pointer to a Lexer
// where the Lexer's io.Reader comes from the named file.
func NewFromFileName(fileName string) *Lexer {
	fd, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Lexer opening file %q for read: %s\n", fileName, err)
	}
	z := NewFromFile(fd)
	z.fileName = fileName
	return z
}

// Next actually calls lexer.nextToken() if it needs to,
// rather than having Consume() actually fetch the next token.
// This allows higher level code to call lexer.Consume() and
// not have it hang if it reads from stdin or something.
func (p *Lexer) Next() (string, TokenType) {
	if p.needsRefresh {
		p.currentToken, p.currentType = p.nextToken()
		p.needsRefresh = false
	}
	return p.currentToken, p.currentType
}

// Consume called by instances of Parse to communicate that
// Parser has used the current token, and will call for the next
// token shortly.
// Just set a flag to actually refresh on subsequent call to Next():
// prevent possible hangs. Consume() doesn't really need to fetch
// the next token, so if it's not available (pipe or stdin), this
// doesn't cause an apparent hang.
func (p *Lexer) Consume() {
	p.needsRefresh = true
}

func (p *Lexer) scan() bool {
	if !p.scanner.Scan() {
		err := p.scanner.Err()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Lexer reading %s: %s\n", p.fileName, err)
		}
		return false
	}

	return true
}

func (p *Lexer) nextToken() (string, TokenType) {

	if worked := p.scan(); !worked {
		return "", EOF
	}

	token := p.scanner.Text()

	// This is kind of bunk, as plSplitter() knows perfectly well what
	// type the token had, but unless I use a package-level variable,
	// I can't figure out how to communicate token type from plSplitter()
	// through bufio.Scanner
	var typ TokenType
	switch token {
	case "~":
		typ = NOT
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

// TokenName returns a human-understandable string
// text reprsentation of the TokenType value you give it.
func TokenName(t TokenType) string {
	r := "unknown"
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
		end := advance + w

		switch c {
		case '(', ')', '&', '~', '|', '=', '>':
			if len(token) == 0 {
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
				token = append(token, data[advance:end]...)
				advance = end
			}
			foundToken = true
		default:
			if c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
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

// BinaryOperator returns true if you pass it one of the
// binary infix propositional logic connectives.
func BinaryOperator(t TokenType) bool {
	switch t {
	case AND, OR, IMPLIES, EQUIV:
		return true
	}
	return false
}

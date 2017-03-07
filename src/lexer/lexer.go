package lexer

import (
	"bufio"
	"log"
	"os"
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
	return "dork", IDENT
}

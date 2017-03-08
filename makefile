all: tokentest parsetest

tokentest: tokentest.go src/lexer/lexer.go
	go build tokentest.go

parsetest: parsetest.go src/lexer/lexer.go src/parser/parser.go src/node/node.go
	go build parsetest.go


clean:
	-rm -rf tokentest parsetest

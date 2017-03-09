all: tokentest parsetest recognizer

tokentest: tokentest.go src/lexer/lexer.go
	go build tokentest.go

recognizer: recognizer.go src/lexer/lexer.go src/parser/parser.go src/node/node.go src/parser/recognizer.go
	go build recognizer.go

parsetest: parsetest.go src/lexer/lexer.go src/parser/parser.go src/node/node.go
	go build parsetest.go


clean:
	-rm -rf tokentest parsetest recognizer
	-rm -rf test_output

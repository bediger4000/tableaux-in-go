all: truthtable

tokentest: tokentest.go src/lexer/lexer.go src/stringbuffer/stringbuffer.go
	go build tokentest.go

recognizer: recognizer.go src/lexer/lexer.go src/parser/parser.go src/node/node.go src/parser/recognizer.go
	go build recognizer.go

parsetest: parsetest.go src/lexer/lexer.go src/parser/parser.go src/node/node.go
	go build parsetest.go

truthtable: truthtable.go src/lexer/lexer.go src/parser/parser.go src/node/node.go src/stringbuffer/stringbuffer.go
	go build truthtable.go


clean:
	-rm -rf tokentest parsetest recognizer truthtable
	-rm -rf test_output

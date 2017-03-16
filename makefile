all: truthtable

tokentest: tokentest.go src/lexer/lexer.go
	go build tokentest.go

recognizer: recognizer.go src/lexer/lexer.go src/parser/parser.go src/node/node.go src/parser/recognizer.go
	go build recognizer.go

parsetest: parsetest.go src/lexer/lexer.go src/parser/parser.go src/node/node.go
	go build parsetest.go

truthtable: truthtable.go src/lexer/lexer.go src/parser/parser.go src/node/node.go 
	go build truthtable.go

tableaux: tableaux.go src/lexer/lexer.go src/parser/parser.go src/node/node.go \
	src/tableaux/tnode.go
	go build tableaux.go


clean:
	-rm -rf tokentest parsetest recognizer truthtable tableaux
	-rm -rf test_output

all: truthtable tableaux

tokentest: tokentest.go src/lexer/lexer.go
	go build tokentest.go

recognizer: recognizer.go src/lexer/lexer.go src/parser/parser.go src/node/node.go src/parser/recognizer.go
	go build recognizer.go

parsetest: parsetest.go src/lexer/lexer.go src/parser/parser.go src/node/node.go
	go build parsetest.go

truthtable: truthtable.go src/lexer/lexer.go src/parser/parser.go src/node/node.go 
	export GOPATH=$$PWD; go build truthtable.go

tableaux: tableaux.go src/lexer/lexer.go src/parser/parser.go src/node/node.go \
	src/tableaux/tnode.go
	export GOPATH=$$PWD; go build tableaux.go

# Need to have GraphViz installed for this to work.
diagrams: tableaux parsetest
	./parsetest -g examplep.dot '~(p&q)=(~p|~q)'
	./tableaux -g examplet.dot '~(p&q)=(~p|~q)'
	dot -Tpng -o examplep.png examplep.dot
	dot -Tpng -o examplet.png examplet.dot

clean:
	-rm -rf tokentest parsetest recognizer truthtable tableaux
	-rm -rf test_output
	-rm -rf *.dot

pl: tokentest.go src/lexer/lexer.go
	go build tokentest.go

clean:
	-rm -rf tokentest

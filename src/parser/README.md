# Parser

Recursive descent parser, based on a
[recursive descent parser for algebraic expressions](https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm)

## Grammar

    EQUIVALENCE -> IMPLICATION {"=" IMPLICATION}
    IMPLICATION -> DISJUNCTION {">" DISJUNCTION}
    DISJUNCTION -> CONJUNCTION {"|" CONJUNCTION}
    CONJUNCTION -> FACTOR {"&" FACTOR}
    FACTOR -> identifier | "(" EQUIVALENCE ")" | "~" FACTOR

The `{something somethingelse}` notation means "a sequence of these types of tokens".

## Recognizer Grammar

    E -> P {BINARYOP P}
    P -> identifier | "(" E ")" | "~" P
    BINARYOP -> "&" | "|" | ">" | "="

The Recognizer Grammar is quite a bit simpler, so I did it first to get my toes wet,
and debug `lexer` methods and functions, and `parser` utility functions. This leads
directly to the recognizer.go program.

## Notes

It seems I broke things up differently than a traditional C `lex` and `yacc` combo.
I put the token types in `package lexer`, where you'd have `yacc` generate them
into `y.tab.h`, which generated lexer `lex.yy.c` would include.

The grammar above led to parser methods like this:

    func (p *Parser) parseEquivalence() (*node.Node) {
        n := p.parseImplication()
        if n != nil {
            for _, typ := p.lexer.Next(); typ == lexer.EQUIV; _, typ = p.lexer.Next() {
                p.lexer.Consume()
                tmp := node.NewOpNode(lexer.EQUIV)
                tmp.Left = n
                tmp.Right = p.parseImplication()
                n = tmp
            }
        }
        return n
    }

And nearly identical methods like this:

    func (p *Parser) parseImplication() (*node.Node) {
        n := p.parseDisjunction()
        if n != nil {
            for _, typ := p.lexer.Next(); typ == lexer.IMPLIES; _, typ = p.lexer.Next() {
                p.lexer.Consume()
                tmp := node.NewOpNode(lexer.IMPLIES)
                tmp.Left = n
                tmp.Right = p.parseDisjunction()
                n = tmp
            }
        }
        return n
    }

The only differences between the two above: `parseImplication` vs `parseDisjunction`
and `EQUIV` vs `IMPLIES`. I can factor this out.

I combined parseEquivalence(), parseImplication(), parseConjunction()
and parseDisjunction() into a single "generalized" function which takes
a `lexer.TokenType` argument. That argument is used to decide what the
condition in the for-loop stops on, and to choose the next function to
call, `parseProdction()` or `parseFactor()`. This generalized function
is a lot harder to understand, but writing `fn := p.parseFactor`, and having
it work is just too cool.

    func (p *Parser) parseProduction(op lexer.TokenType ) (*node.Node) {
    
        nextProduction := p.parseProduction
        if op == lexer.AND {
            nextProduction = p.parseFactor
        }
    
        no := nextOp[op]
        newNode := nextProduction(no)
        if newNode != nil {
            for _, typ := p.lexer.Next(); typ == op; _, typ = p.lexer.Next() {
                p.lexer.Consume()
                tmp := node.NewOpNode(op)
                tmp.Left = newNode
                tmp.Right = nextProduction(no)
                newNode = tmp
            }
        }
        return newNode
    }

The keys in this method are `no := nextOp[op]`, which is the operation getting factored out,
and `nextProduction := p.parseProduction`. Not unlike Python, you can keep a reference to a
particular object and one of its methods.

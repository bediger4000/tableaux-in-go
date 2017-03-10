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
and debug `lexer` methods and functions, and `parser` utility functions

## Notes

It seems I broke things up differently than a traditional C `lex` and `yacc` combo.
I put the token types in `package lexer`, where you'd have `yacc` generate them
into `y.tab.h`, which generated lexer `lex.yy.c` would include.

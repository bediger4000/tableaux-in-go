# Prove propositional logic tautologies via Smullyan's analytic tableaux method 

## Something

## Something else!

## Software engineering notes

I started with a Golang [parser for propositional logic expressions](https://github.com/bediger4000/propositional-logic-go),
v2.0 . The idea was to edge up to a tableaux method tautology prover.

0. Write lexer for propositional logic tokens
1. Write propositional logic parser
2. Write truth table generator based on the parser
3. Write single-level tableaux generator
4. Write full-on tableaux prover 

## Other programs in this project

    ./truthtable '(p&q>...)'

Prints a truth table for the command line expression. The expression gets parsed and
evaluated with every combination of true and false for each variable. Can be helpful
verifying whether `tableaux` gets its proof correct.

     ./tokentest 'a&b&c())~|>='

Runs the lexer over the command line argument, and prints out all the tokens it finds,
and their types. Used to write and debug package `lexer`

     ./recognizer test_input/001

Uses a simpler recursive descent grammar than `tableaux` does to recognize propositional logic
formula. Used to learn how to write a recursive descent grammar in Go.

    ./parsetest test_input/004

Lexes, parses, then prints the propositional logic expressions in the file named
on command line. Used to develop and debug package `parser`

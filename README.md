# Prove propositional logic tautologies via Smullyan's analytic tableaux method 

Famous logician [Raymond Smullyan ](https://en.wikipedia.org/wiki/Raymond_Smullyan)
used a proof procedure called [analytic tableaux](https://en.wikipedia.org/wiki/Method_of_analytic_tableaux)
to prove tautologies in propositional logic.

This program is based on chapters from three books by Smullyan:

* _A Beginner's Guide to Mathematical Logic_, Dover, 2014, chapter 6
* _Logical Labyrinths_, CRC Press, 2009, chapter 11
* _First Order Logic_, Dover, 1995, chapter II

All these books have essentially the same explanation with slight variations.
This project does signed tableaux.

`tableaux` supports these binary infix logical oparators:

* `&` - conjunction
* `|` - disjunction
* `>` - material implication
* `=` - logical equivalence

And one binary prefix operator for negation: `~`

## Building the program

    $ make tableaux

## Using the program

Invoked with a single propositional logic expression, `tableaux`
writes out a tableau that proves whether the expression constitutes
a tautology or not.

    $ ./tableaux '((p>q)>r) > ((p>q)>(p>r))'
    Expression: "((p > q) > r) > ((p > q) > (p > r))"
    /*

    0. false: ((p > q) > r) > ((p > q) > (p > r))
    1. true: (p > q) > r (0)
    2. false: (p > q) > (p > r) (0)
       3 left, 4 right

    3. false: p > q (1)
    5. true: p > q (2) contradicts 3


    4. true: r (1)
    6. true: p > q (2)
    7. false: p > r (2)
       8 left, 9 right

    8. false: p (6)
    10. true: p (7) contradicts 8


    9. true: q (6)
    11. true: p (7)
    12. false: r (7) contradicts 4

    Formula is a tautology
    */

Called with more than one propositional logic expression, `tableaux` proves
whether or not the final expression is a logical consequence of the other expressions.

## Proof Procedure

As pseudocode:

    do {
        find all unclosed leaf nodes of tableau

        if no unclosed leave nodes exist:
            the expression is tautological

        for each unclosed leaf node:

            Find an unused forumla as far up the tableaux as possible
            on the branch that the unclosed leaf node resides on.

            if such an unused formula exits:

                Subjoin inferences of the unused formula to all
                unclosed leaf nodes beneath it currently in the tableaux.
                Mark inferences that consist of a signed identifer as used.

                Check each newly-subjoined inference for contradictions with
                previous inferences in the tableau's branch it resides on. Mark
                any leaf nodes that cause a contradiction as closed.

                Mark the unused formula as used.

                exit for-each loop over unclosed leaf nodes.

    } while an unused formula was found

This algorithm terminates, since each formula gets used to subjoin inferences
only once. Subjoined inferences that contradict previously subjoined inferences,
"close" a branch of the tableaux so that no further inferences are added to that branch. 

## Parsing and Lexing

See [README for package parser](https://github.com/bediger4000/tableaux-in-go/tree/master/src/parser)  for details on this topic.

### Parse Tree

Note that a parse tree constructed by packages `lexer` and `parser` differs
from a (finished) tableau.


![Parse tree for ~(p&q)=(~p|~q)](https://raw.githubusercontent.com/bediger4000/tableaux-in-go/master/examplep.png)

*Parse tree for ~(p&q)=(~p|~q)*

![Finished tableau for`~(p&q)=(~p|~q)`](https://raw.githubusercontent.com/bediger4000/tableaux-in-go/master/examplet.png)

*Finished tableau for ~(p&q)=(~p|~q)*

## Data structure for tableaux

Incautious reading any of the 3 Smullyan books  above on analytic tableaux would have you
believe that an analytic tableaux consists of arrays of subexpressions of the propositional
logic formula to be proved. Smullyan was not a programmer, it seems, because a tableaux is
a binary tree of individual subexpressions. The sign part of an analytic tableaux is attached
to a subexpression, as is the notion of an unused formula, and whether a branch is closed or not.

    type Tnode struct {
        Sign       bool
        Tree       *node.Node
        Expression string
    
        Parent     *Tnode
        Left       *Tnode
        Right      *Tnode
    
        Used       bool
        closed     bool
    }

The `Sign` and `Expression` members of this struct identify a node in a tableau. Apparent duplicate 
nodes can appear in a tableau, because the logic-not handling for a program is more literal than
for a human. Because the proof procedure checks for contractions with previous subexpressions in
a tableau branch, a node has to have a link back to its "parent" in the tableau. Only the root
node of a tableau's binary tree has a nil value for `Parent`. Each node in a tableau keeps a pointer
to the node of a parse tree that corresponds to the tableau node itself. Subjoining inferences
to leaf nodes of a branch uses the principal connective of its pare tree pointer to decide
how to subjoin (linearly or bifurcate), and the sign of the subjoined expressions.

Typing `Tnode.Sign` as a Golang boolean is semantically obvious: the signs of expressions
in Smullyan's tableaux are 'T' or 'F', but internally, a program could use 0 and 1, or even
two different strings. Checking two lines in a tableau (two nodes in a binary tree) for
contradiction only involves non-equality of the `Sign` element, and string equality of the
`Expression` element. A program could check the `Tree` elements for tree-equality, but since
`tableaux.Print()` canonicalizes string representations of tableau binary trees, string equality
is sufficient.

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

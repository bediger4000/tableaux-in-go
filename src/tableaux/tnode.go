// Smullyan's Analytic Tableaux, as a Go type.
package tableaux

// See:
// "A Beginner's Guide to Mathematical Logic", Dover, 2014, chapter 6
// "Logical Labyrinths", CRC Press, 2009, chapter 11
// "First Order Logic", Dover, xxxx, chapter N
// for essentially the same explanation with slight variations.
// This does signed tableaux.

import (
	"fmt"
	"io"
	"lexer"
	"node"
)

type Tnode struct {
	// Set in New(), should never get changed
	Sign       bool
	Tree       *node.Node
	Expression string  // element Tree as a string.

	// Changed during subjoining inferences, and initial setup.
	Parent     *Tnode
	Left       *Tnode
	Right      *Tnode

	Used       bool    // Have interence(s) of this expression been subjoined to leaf nodes?
	closed     bool    // Does this expression contradict a predecessor in the tableau?
}

// The only way to create a Tnode instance.
func New(tree *node.Node, sign bool, parent *Tnode) (*Tnode) {
	var r Tnode

	r.Tree = tree
	r.Used = false
	if tree.Op == lexer.IDENT {
		r.Used = true  // No inferences to make from an identifier.
	}
	r.closed = false
	r.Parent = parent
	r.Sign = sign
	r.Expression = node.ExpressionToString(r.Tree)

	return &r
}

// Find all unclosed leaf node(s) below the receiver in
// a tableau. Leaf might be marked "used" if it's just an identifier,
// also this can return zero-len array if all leaf nodes marked closed
func (n *Tnode) FindUnclosedLeaf() ([]*Tnode) {
	var a []*Tnode
	if n.Left == nil && n.Right == nil {
		if !n.closed {
			a = append(a, n)
		}
	}
	if n.Left != nil {
		t := n.Left.FindUnclosedLeaf()
		a = append(a, t...)
	}
	if n.Right != nil {
		t := n.Right.FindUnclosedLeaf()
		a = append(a, t...)
	}
	return a
}

// Find an unused formula above an unclosed leaf node
// by following the Tnode.Parent links up a branch.
func (n *Tnode) FindTallestUnused() *Tnode {
	var p *Tnode
	var unused *Tnode

	// Have to consider n (the unclosed leaf node itself)
	// as it might be the only unused Tnode in the branch.
	// Also have to walk Tnode.Parent chain all the way up
	// to the root of the tableau, because IDENT node.Node
	// objects can appear below an unused node.Node in a branch.
	for p = n; p != nil; p = p.Parent {
		if !p.Used {
			unused = p
		}
	}
	return unused
}

// Try to find a contradiction to receiver Tnode instance n by following
// Tnode.Parent links all the way up a branch of a tableau
// Not recursive, so the receiver n is the expression possibly contradicted by
// element further back up the tableau branch.
func (n *Tnode) CheckForContradictions() bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if n.Sign != p.Sign && n.Expression == p.Expression {
			fmt.Printf("Leaf %v: %q contradicted by ancestor %v: %q\n", n.Sign, n.Expression, p.Sign, p.Expression)
			n.closed = true
			return true
		}
	}
	return false
}

// Subjoin inferences of from to Tnode instance named parent.
func (parent *Tnode) AddInferences(from *Tnode) {

	if from.Tree.Op == lexer.IDENT {
		return
	}

	// Smullyan's beta-type, and logical equivalence. These create
	// bifurcations in branches, so do them first. Alpha-type inferences,
	// which just linearly extend a branch, push inferences down the branch.
	// By doing beta-type, bifurcating inferences first, branches are made
	// linearly as long as possible.

	if (from.Tree.Op == lexer.AND && from.Sign == false) || (from.Tree.Op == lexer.OR && from.Sign == true) {
		immediate := New(from.Tree.Left, from.Sign, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q left of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		immediate.CheckForContradictions()

		immediate2 := New(from.Tree.Right, from.Sign, parent)
		fmt.Printf("Adding %v: %q right of %v: %q\n", immediate2.Sign, immediate2.Expression, parent.Sign, parent.Expression)
		parent.Right = immediate2

		immediate2.CheckForContradictions()

		return
	}

	if from.Tree.Op == lexer.IMPLIES && from.Sign == true {
		immediate := New(from.Tree.Left, false, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q left of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		immediate.CheckForContradictions()

		immediate2 := New(from.Tree.Right, true, parent)
		parent.Right = immediate2
		fmt.Printf("Adding %v: %q right of %v: %q\n", immediate2.Sign, immediate2.Expression, parent.Sign, parent.Expression)

		immediate2.CheckForContradictions()

		return
	}

	// Not actually a beta-type, and Smullyan probably would seems rather
	// define equivalance as an abbreviation. It does create a new bifurcation
	// in a branch, however.
	if from.Tree.Op == lexer.EQUIV {

		var sign1, sign2, sign3, sign4 bool
		if from.Sign == true {
			sign1, sign2, sign3, sign4 = true, true, false, false
		} else {
			sign1, sign2, sign3, sign4 = true, false, false, true
		}

		immediate1 := New(from.Tree.Left, sign1, parent)
		parent.Left = immediate1
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate1.Sign, immediate1.Expression, parent.Sign, parent.Expression)

		if !immediate1.CheckForContradictions() {

			immediate2 := New(from.Tree.Right, sign2, immediate1)
			immediate1.Left = immediate2
			fmt.Printf("Adding %v: %q below of %v: %q\n", immediate2.Sign, immediate2.Expression, immediate1.Sign, immediate1.Expression)

			immediate2.CheckForContradictions()
		}

		immediate3 := New(from.Tree.Left, sign3, parent)
		parent.Right = immediate3
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate3.Sign, immediate3.Expression, parent.Sign, parent.Expression)

		if !immediate3.CheckForContradictions() {

			immediate4 := New(from.Tree.Right, sign4, immediate3)
			immediate3.Left = immediate4
			fmt.Printf("Adding %v: %q below of %v: %q\n", immediate4.Sign, immediate4.Expression, immediate3.Sign, immediate3.Expression)

			immediate4.CheckForContradictions()
		}

		return
	}

	// Smullyan's alpha-type inferences. These just extend a branch, without bifurcating it.
	// For AND, OR, IMPLIES alpha-type inferences, add the 2nd of two inferences immediately
	// below the parent. Combined with doing beta-type, bifurcating inferencese first, this
	// keeps the branch linear for as long as possible.

	if from.Tree.Op == lexer.NOT {
		immediate := New(from.Tree.Left, !from.Sign, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		parent.Left.CheckForContradictions()
		return
	}

	if (from.Tree.Op == lexer.AND && from.Sign == true) || (from.Tree.Op == lexer.OR && from.Sign == false) {
		immediate := New(from.Tree.Left, from.Sign, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		// Check 1st inference for contradictions, don't bother subjoining 2nd inference
		// if 1st one has a contradction and closes the branch.
		if !immediate.CheckForContradictions() {

			immediate2 := New(from.Tree.Right, from.Sign, immediate)
			immediate.Left = immediate2
			fmt.Printf("Adding %v: %q below of %v: %q\n", immediate2.Sign, immediate2.Expression, immediate.Sign, immediate.Expression)

			immediate2.CheckForContradictions()
		}

		return
	}

	// Material implication causes a special case: F: p>q means that T:p and F:q get subjoined,
	// preventing the general alpha-type code above from working.
	if from.Tree.Op == lexer.IMPLIES && from.Sign == false {
		parent.Left = New(from.Tree.Left, true, parent)
		fmt.Printf("Adding %v: %q below of %v: %q\n", parent.Left.Sign, parent.Left.Expression, parent.Sign, parent.Expression)
		if ! parent.Left.CheckForContradictions() {

			parent.Left.Left = New(from.Tree.Right, false, parent.Left)
			fmt.Printf("Adding %v: %q below of %v: %q\n", parent.Left.Left.Sign, parent.Left.Left.Expression, parent.Left.Sign, parent.Left.Expression)
			parent.Left.Left.CheckForContradictions()
		}

		return
	}

	// Don't think it should ever get here.
	errString := fmt.Sprintf("Trying to add inferences of %v:%q to leaf node %v:%q\n", from.Sign, from.Expression, parent.Sign, parent.Expression)
	panic(errString)
}

// Do the actual work of writing GraphViz digraph output to io.Writer w.
// Another traverse of a tableau (binary tree of *Tnode instances),
// with semantic irregularities causing some inorder and some postorder
// operations.
// The Tnode.Parent backlink can help in debugging.
func (p *Tnode) graphTnode(w io.Writer) {
	sign := "F"
	if p.Sign { sign = "T" }

	// Append a string to the formula, inlucde 'U' for a formula
	// whose inferences got subjoined to all it's leaf nodes,
	// and 'C' for the leaf node of a closed branch.
	extra := ""
	if p.Used {
		extra += "U"
	}
	if p.closed {
		extra += "C"
	}

	fmt.Fprintf(w, "n%p [label=\"%s: %s%s\"];\n", p, sign, p.Expression, ", "+extra)
/*
	if p.Parent != nil {
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Parent)
	}
*/
	if p.Left != nil {
		p.Left.graphTnode(w)
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Left)
	}
	if p.Right != nil {
		p.Right.graphTnode(w)
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Right)
	}
}

// Write GraphViz directed graph dot input to argument w io.Writer.
func (p *Tnode) GraphTnode(w io.Writer) {
	fmt.Fprintf(w, "digraph g {\n")
	p.graphTnode(w)
	fmt.Fprintf(w, "}\n")
}

// Append argument n *Tnode to the leaf node of receiver p in a branch of a tableau.
// This assumes that there's just a linked list via Tnode.Left elements. Used
// only in setting up the hypotheses for finding consequences of a list of
// formulas, so just followin Tnode.Left works.
func (p *Tnode) AppendLeaf(n *Tnode) {
	var leaf *Tnode
	for t := p; t != nil; t = t.Left {
		leaf = t
	}
	leaf.Left = n
	n.Parent = leaf
}

// Shouldn't this just make Tnode match type Stringer?
func (p *Tnode) PrintTnode() {
	fmt.Printf("Tnode %p\n", p)
	fmt.Printf("\ttree %p\n", p.Tree)
	fmt.Printf("\t%v: %q\n", p.Sign, p.Expression)
	fmt.Printf("\tUsed   %v\n", p.Used)
	fmt.Printf("\tclosed %v\n", p.closed)
	fmt.Printf("\tParent %p\n", p.Parent)
	fmt.Printf("\tLeft   %p\n", p.Left)
	fmt.Printf("\tRight  %p\n", p.Right)

	if p.Left != nil {
		p.Left.PrintTnode()
	}
	if p.Right != nil {
		p.Right.PrintTnode()
	}
}

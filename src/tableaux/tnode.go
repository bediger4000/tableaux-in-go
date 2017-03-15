package tableaux

import (
	"fmt"
	"io"
	"lexer"
	"node"
)

type Tnode struct {
	Sign       bool
	tree       *node.Node
	Expression string
	Parent     *Tnode
	Left       *Tnode
	Right      *Tnode
	Used       bool
	closed     bool
}

func New(tree *node.Node, sign bool, parent *Tnode) (*Tnode) {
	var r Tnode

	r.tree = tree
	r.Used = false
	if tree.Op == lexer.IDENT {
		r.Used = true
	}
	r.closed = false
	r.Parent = parent
	r.Sign = sign
	r.Expression = node.ExpressionToString(r.tree)

	return &r
}

// Returns nil if it can't find an unused expression
func (n *Tnode)FirstUnused() (*Tnode) {

	if n == nil {
		return nil
	}

	if !n.Used {
		return n
	}
	var r *Tnode
	r = n.Left.FirstUnused()
	if r == nil {
		r = n.Right.FirstUnused()
	}
	return r
}

// Find an unclosed leaf node - it might
// be marked "used" if it's just an identifier,
// also can return nil if all leaf nodes marked closed
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

// Find an unused formala above an unclosed leaf node
func (n *Tnode) FindTallestUnused() *Tnode {
	// Walk linked list formed by Tnode.Parent pointers
	var p *Tnode
	var unused *Tnode
	for p = n; p != nil; p = p.Parent {
		if !p.Used {
			unused = p
		}
	}
	return unused
}

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

func (parent *Tnode) AddInferences(from *Tnode) {

	if from.tree.Op == lexer.IDENT {
		return
	}

	// Smullyan's beta-type

	if (from.tree.Op == lexer.AND && from.Sign == false) || (from.tree.Op == lexer.OR && from.Sign == true) {
		immediate := New(from.tree.Left, from.Sign, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q left of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		immediate.CheckForContradictions()

		immediate2 := New(from.tree.Right, from.Sign, parent)
		fmt.Printf("Adding %v: %q right of %v: %q\n", immediate2.Sign, immediate2.Expression, parent.Sign, parent.Expression)
		parent.Right = immediate2

		immediate2.CheckForContradictions()

		return
	}

	if from.tree.Op == lexer.IMPLIES && from.Sign == true {
		immediate := New(from.tree.Left, false, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q left of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		immediate.CheckForContradictions()

		immediate2 := New(from.tree.Right, true, parent)
		parent.Right = immediate2
		fmt.Printf("Adding %v: %q right of %v: %q\n", immediate2.Sign, immediate2.Expression, parent.Sign, parent.Expression)

		immediate2.CheckForContradictions()

		return
	}

	// Not actually a beta-type, and Smullyan probably would seems rather
	// define equivalance as an abbreviation.
	if from.tree.Op == lexer.EQUIV {

		var sign1, sign2, sign3, sign4 bool
		if from.Sign == true {
			sign1, sign2, sign3, sign4 = true, true, false, false
		} else {
			sign1, sign2, sign3, sign4 = true, false, false, true
		}

		immediate1 := New(from.tree.Left, sign1, parent)
		parent.Left = immediate1
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate1.Sign, immediate1.Expression, parent.Sign, parent.Expression)

		if !immediate1.CheckForContradictions() {

			immediate2 := New(from.tree.Right, sign2, immediate1)
			immediate1.Left = immediate2
			fmt.Printf("Adding %v: %q below of %v: %q\n", immediate2.Sign, immediate2.Expression, immediate1.Sign, immediate1.Expression)

			immediate2.CheckForContradictions()
		}

		immediate3 := New(from.tree.Left, sign3, parent)
		parent.Right = immediate3
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate3.Sign, immediate3.Expression, parent.Sign, parent.Expression)

		if !immediate3.CheckForContradictions() {

			immediate4 := New(from.tree.Right, sign4, immediate3)
			immediate3.Left = immediate4
			fmt.Printf("Adding %v: %q below of %v: %q\n", immediate4.Sign, immediate4.Expression, immediate3.Sign, immediate3.Expression)

			immediate4.CheckForContradictions()
		}

		return
	}

	// Smullyan's alpha-type
	if from.tree.Op == lexer.NOT {
		immediate := New(from.tree.Left, !from.Sign, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		parent.Left.CheckForContradictions()
		return
	}

	if (from.tree.Op == lexer.AND && from.Sign == true) || (from.tree.Op == lexer.OR && from.Sign == false) {
		immediate := New(from.tree.Left, from.Sign, parent)
		parent.Left = immediate
		fmt.Printf("Adding %v: %q below of %v: %q\n", immediate.Sign, immediate.Expression, parent.Sign, parent.Expression)

		if !immediate.CheckForContradictions() {

			immediate2 := New(from.tree.Right, from.Sign, immediate)
			immediate.Left = immediate2
			fmt.Printf("Adding %v: %q below of %v: %q\n", immediate2.Sign, immediate2.Expression, immediate.Sign, immediate.Expression)

			immediate2.CheckForContradictions()
		}

		return
	}

	if from.tree.Op == lexer.IMPLIES && from.Sign == false {
		parent.Left = New(from.tree.Left, true, parent)
		fmt.Printf("Adding %v: %q below of %v: %q\n", parent.Left.Sign, parent.Left.Expression, parent.Sign, parent.Expression)
		if ! parent.Left.CheckForContradictions() {

			parent.Left.Left = New(from.tree.Right, false, parent.Left)
			fmt.Printf("Adding %v: %q below of %v: %q\n", parent.Left.Left.Sign, parent.Left.Left.Expression, parent.Left.Sign, parent.Left.Expression)
			parent.Left.Left.CheckForContradictions()
		}

		return
	}
}

func (p *Tnode) graphTnode(w io.Writer) {
	sign := "F"
	if p.Sign { sign = "T" }
	extra := ""
	if p.Used {
		extra += " U"
	}
	if p.closed {
		extra += " C"
	}
	fmt.Fprintf(w, "n%p [label=\"%s: %s%s\"];\n", p, sign, p.Expression, extra)
	if p.Left != nil {
		p.Left.graphTnode(w)
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Left)
	}
	if p.Right != nil {
		p.Right.graphTnode(w)
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Right)
	}
}

func (p *Tnode) GraphTnode(w io.Writer) {
	fmt.Fprintf(w, "digraph g {\n")
	p.graphTnode(w)
	fmt.Fprintf(w, "}\n")
}

func (p *Tnode) PrintTnode() {
	fmt.Printf("Tnode %p\n", p)
	fmt.Printf("\ttree %p\n", p.tree)
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

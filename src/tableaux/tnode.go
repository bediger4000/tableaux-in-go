package tableaux

import (
	"fmt"
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
func FindUnclosedLeaf(n *Tnode) ([]*Tnode) {
	var a []*Tnode
	if n.Left == nil && n.Right == nil {
		if !n.closed {
			a = append(a, n)
		}
	}
	if n.Left != nil {
		t := FindUnclosedLeaf(n.Left)
		a = append(a, t...)
	}
	if n.Right != nil {
		t := FindUnclosedLeaf(n.Right)
		a = append(a, t...)
	}
	return a
}

// Find an unused formala above an unclosed leaf node
func (n *Tnode) FindTallestUnused() *Tnode {
	// Walk linked list formed by Tnode.Parent pointers
	var p *Tnode
	for p = n.Parent; p != nil; p = p.Parent {
		if p.Used {
			break
		}
	}
	// If p == nil, got to root w/o finding an unused formula
	return p
}

func (n *Tnode) CheckForContradictions() {
	if n.Left != nil {
		n.Left.CheckForContradictions()
	}
	if n.Left == nil && n.Right == nil {
		if !n.closed {
			for p := n.Parent; p.Parent != nil; p = p.Parent {
				if n.Sign != p.Sign && n.Expression == p.Expression {
					n.closed = true
				}
			}
		}
	}
	if n.Right != nil {
		n.Right.CheckForContradictions()
	}
}

func (n *Tnode) SubjoinInferences(unused *Tnode) {
	if n.Left == nil && n.Right == nil {
		if !n.closed {
			n.AddInferences(unused)
		}
		return
	}
	if n.Left != nil {
		n.Left.SubjoinInferences(unused)
	}
	if n.Right != nil {
		n.Right.SubjoinInferences(unused)
	}
}

func (parent *Tnode) AddInferences(p *Tnode) {

	if p.tree.Op == lexer.IDENT {
		return
	}

	// Smullyan's alpha-type first
	if p.tree.Op == lexer.NOT {
		immediate := New(p.tree.Left, !p.Sign, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate
		return
	}

	if (p.tree.Op == lexer.AND && p.Sign == true) || (p.tree.Op == lexer.OR && p.Sign == false) {
		immediate := New(p.tree.Left, p.Sign, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate
		immediate2 := New(p.tree.Right, p.Sign, immediate)
		if immediate.Left != nil {
			tmp := immediate.Left
			tmp.Parent = immediate2
			immediate2.Left = tmp
		}
		immediate.Left = immediate2
		return
	}

	if p.tree.Op == lexer.IMPLIES && p.Sign == false {
		immediate := New(p.tree.Left, true, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate
		immediate2 := New(p.tree.Right, false, immediate)
		if immediate.Left != nil {
			tmp := immediate.Left
			tmp.Parent = immediate2
			immediate2.Left = tmp
		}
		immediate.Left = immediate2
		return
	}

	// Smullyan's beta-type
	if (p.tree.Op == lexer.AND && p.Sign == false) || (p.tree.Op == lexer.OR && p.Sign == true) {
		immediate := New(p.tree.Left, p.Sign, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate

		immediate2 := New(p.tree.Right, p.Sign, parent)
		p.Right = immediate2

		return
	}

	if p.tree.Op == lexer.IMPLIES && p.Sign == true {
		immediate := New(p.tree.Left, false, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate

		immediate2 := New(p.tree.Right, true, parent)
		p.Right = immediate2

		return
	}
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

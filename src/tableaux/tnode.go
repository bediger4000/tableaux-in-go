package tableaux

import (
	"fmt"
	"lexer"
	"node"
)

type Tnode struct {
	sign bool
	tree         *node.Node
	expression string
	Parent       *Tnode
	Left         *Tnode
	Right        *Tnode
	used         bool
	closed       bool
}

func New(tree *node.Node, sign bool, parent *Tnode) (*Tnode) {
	var r Tnode

	r.tree = tree
	r.used = false
	r.closed = false
	r.Parent = parent
	r.sign = sign
	r.expression = node.ExpressionToString(r.tree)

	return &r
}

// Returns nil if it can't find an unused expression
func FirstUnused() (n *Tnode) {
	if n.used {
		return n
	}
	var r *Tnode
	if n.Left {
		r = FirstUnused(n.Left)
	}
	if r == nil {
		r = FirstUnused(n.Right)
	}
	return r
}

func (n *Tnode) Subjoin(unused *Tnode) {
	if n.Left == nil && n.Right == nil {
		if !n.closed {
		}
		return
	}
	if n.Left != nil {
		n.Left.Subjoin(unused)
	}
	if n.Right != nil {
		n.Right.Subjoin(unused)
	}
}

func (p *Tnode) DoTnode(parent *Tnode) {

	if p.tree.Op == lexer.IDENT {
		p.used = true
		return
	}

	// Smullyan's alpha-type first
	if p.tree.Op == lexer.NOT {
		immediate := New(p.tree.Left, !p.sign, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate
		p.used = true
		return
	}

	if (p.tree.Op == lexer.AND && p.sign == true) || (p.tree.Op == lexer.OR && p.sign == false) {
		immediate := New(p.tree.Left, p.sign, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate
		immediate2 := New(p.tree.Right, p.sign, immediate)
		if immediate.Left != nil {
			tmp := immediate.Left
			tmp.Parent = immediate2
			immediate2.Left = tmp
		}
		immediate.Left = immediate2
		p.used  = true
		return
	}

	if p.tree.Op == lexer.IMPLIES && p.sign == false {
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
		p.used  = true
		return
	}

	// Smullyan's beta-type
	if (p.tree.Op == lexer.AND && p.sign == false) || (p.tree.Op == lexer.OR && p.sign == true) {
		immediate := New(p.tree.Left, p.sign, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate

		immediate2 := New(p.tree.Right, p.sign, parent)
		p.Right = immediate2

		p.used  = true
		return
	}

	if p.tree.Op == lexer.IMPLIES && p.sign == true {
		immediate := New(p.tree.Left, false, parent)
		if p.Left != nil {
			tmp := p.Left
			tmp.Parent = immediate
			immediate.Left = tmp
		}
		p.Left = immediate

		immediate2 := New(p.tree.Right, true, parent)
		p.Right = immediate2

		p.used  = true
		return
	}
}

func (p *Tnode) PrintTnode() {
	fmt.Printf("Tnode %p\n", p)
	fmt.Printf("\ttree %p\n", p.tree)
	fmt.Printf("\t%v: %q\n", p.sign, p.expression)
	fmt.Printf("\tused   %v\n", p.used)
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

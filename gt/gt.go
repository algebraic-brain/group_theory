package gt

import (
	"fmt"
	"sync"
)

var lastToken = 1
var mut = &sync.Mutex{}

func tok() int {
	mut.Lock()
	token := lastToken
	lastToken++
	mut.Unlock()
	return token
}

//Grooup element interface:
type Element interface {
	//Checks whether two elements are equal literally (although Same() may return false)
	EqualLiteral(Element) bool
	//Makes literal clone of element (although Same() will return false)
	CloneLiteral() Element
	//Checks whether one of two elements was made from other during the proof
	Same(Element) bool

	ToComposite() *Composite
	ToInversed() *Inversed
	ToNamed() *Named
	ToIdentity() *Identity

	setToken(int)
	token() int
	same(Element) bool
}

//Element prototype:
type element struct {
	tok int
}

func (el *element) token() int {
	return el.tok
}

func (el *element) init() {
	el.tok = tok()
}

func (el *element) setToken(t int) {
	el.tok = t
}

func (el *element) ToComposite() *Composite { panic("It's not Composite") }
func (el *element) ToInversed() *Inversed   { panic("It's not Inversed") }
func (el *element) ToNamed() *Named         { panic("It's not Named") }
func (el *element) ToIdentity() *Identity   { panic("It's not Identity") }

//Checks whether one of two elements was made from other during the proof
func (el *element) Same(other Element) bool {
	return el.same(other)
}

func (el *element) same(other Element) bool {
	return el.tok == other.token()
}

//Verify proof (forth, back) that $left = right$
func Verify(left, right Element, forth, back func(Element) Element) bool {
	l := left.CloneLiteral()
	r := right.CloneLiteral()

	lr := forth(l)
	rl := back(r)

	backIsStep := r.same(rl)
	forthIsStep := l.same(lr)

	return forthIsStep && backIsStep && lr.EqualLiteral(r) && rl.EqualLiteral(l)
}

func VerifyForth(left, right Element, forth func(Element) Element) bool {
	l := left.CloneLiteral()
	r := right.CloneLiteral()

	lr := forth(l)
	forthIsStep := l.same(lr)

	if !forthIsStep {
		fmt.Println("VerifyForth: 'forth' is not step")
	}

	lrEq := lr.EqualLiteral(r)

	if !lrEq {
		fmt.Println("VerifyForth: 'forth(left) != right'")
	}
	return forthIsStep && lr.EqualLiteral(r)
}

//Composite element of group:
type Composite struct {
	element
	left  Element
	right Element
}

func (c *Composite) ToComposite() *Composite { return c }

//Composes elements of group
func Compose(a, b Element) *Composite {
	c := &Composite{
		left:  a.CloneLiteral(),
		right: b.CloneLiteral(),
	}
	c.init()
	return c
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Composite) EqualLiteral(other Element) bool {
	if o, ok := other.(*Composite); ok {
		return o.left.EqualLiteral(c.left) && o.right.EqualLiteral(c.right)
	}
	return false
}

//Makes literal clone of element (although same() will return false)
func (c *Composite) CloneLiteral() Element {
	return Compose(c.left.CloneLiteral(), c.right.CloneLiteral())
}

//turns $a\cdot (b\cdot c)$ to $(a\cdot b)\cdot c$. This is a step of proof.
func (c *Composite) Associate() *Composite {
	if r, ok := c.right.(*Composite); ok {
		n := Compose(Compose(c.left, r.left), r.right)
		n.setToken(c.token())
		return n
	}
	panic("Associator requires $a\\cdot (b\\cdot c)$ type arguments")
}

//turns $(a\cdot b)\cdot c$ to $a\cdot (b\cdot c)$. This is a step of proof.
func (c *Composite) Unassociate() *Composite {
	if l, ok := c.left.(*Composite); ok {
		n := Compose(l.left, Compose(l.right, c.right))
		n.setToken(c.token())
		return n
	}
	panic("Unassociator requires $a\\cdot (b\\cdot c)$ type arguments")
}

//turns $a\cdot a^{-1}$ and $a^{-1}\cdot a$ to $e$. This is a step of proof.
func (c *Composite) Annihilate() *Identity {
	n := NewIdentity()
	n.setToken(c.token())
	if r, ok := c.right.(*Inversed); ok {
		if r.operand.EqualLiteral(c.left) {
			return n
		}
	} else if l, ok := c.left.(*Inversed); ok {
		if l.operand.EqualLiteral(c.right) {
			return n
		}
	}
	panic("Annihilator requires $a\\cdot (a^{-1})$ type arguments")
}

//turns $a\cdot e$ and $e\cdot a$ to $a$. This is a step of proof.
func (c *Composite) Simplify() Element {
	if _, ok := c.right.(*Identity); ok {
		n := c.left
		n.setToken(c.token())
		return n
	} else if _, ok := c.left.(*Identity); ok {
		n := c.right
		n.setToken(c.token())
		return n
	}
	panic("Simplificator requires $a\\cdot e$ or  $e\\cdot a$ type arguments")
}

//turns $a$ and $e\cdot a$ or $a\cdot e$ depending on "left". This is a step of proof.
func Unsimplify(el Element, left bool) *Composite {
	if left {
		n := Compose(NewIdentity(), el)
		n.setToken(el.token())
		return n
	} else {
		n := Compose(el, NewIdentity())
		n.setToken(el.token())
		return n
	}
}

//maps proofs to left and right elements of composite. This is a step of proof iff both "left" and "right" are steps.
func (c *Composite) Map(left func(Element) Element, right func(Element) Element) *Composite {
	l := left(c.left)
	r := right(c.right)
	n := Compose(l, r)
	if l.same(c.left) && r.same(c.right) {
		n.setToken(c.token())
	}
	return n
}

//returns left element of composite
func (c *Composite) Left() Element {
	return c.left.CloneLiteral()
}

//returns right element of composite
func (c *Composite) Right() Element {
	return c.right.CloneLiteral()
}

//Inversed element of group:
type Inversed struct {
	element
	operand Element
}

func (c *Inversed) ToInversed() *Inversed { return c }

//Inverses element of group:
func Inverse(el Element) *Inversed {
	n := &Inversed{
		operand: el,
	}
	n.init()
	return n
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Inversed) EqualLiteral(other Element) bool {
	if o, ok := other.(*Inversed); ok {
		return o.operand.EqualLiteral(c.operand)
	}
	return false
}

//Makes literal clone of element (although same() will return false)
func (c *Inversed) CloneLiteral() Element {
	return Inverse(c.operand.CloneLiteral())
}

//maps proofs to operand of inversion. This is a step of proof iff "f" is a step.
func (c *Inversed) Map(f func(Element) Element) *Inversed {
	op := f(c.operand)
	n := Inverse(op)
	if op.same(c.operand) {
		n.setToken(c.token())
	}
	return n
}

//Ordinary named element of group
type Named struct {
	element
	name string
}

func (c *Named) ToNamed() *Named { return c }

//returns name of named element
func (c *Named) Name() string {
	return c.name
}

//Creates new named element
func NewNamed(name string) *Named {
	n := &Named{
		name: name,
	}
	n.init()
	return n
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Named) EqualLiteral(other Element) bool {
	if o, ok := other.(*Named); ok {
		return o.name == c.name
	}
	return false
}

//Makes literal clone of element (although same() will return false)
func (c *Named) CloneLiteral() Element {
	return NewNamed(c.name)
}

//Identity element
type Identity struct {
	element
}

func (c *Identity) ToIdentity() *Identity { return c }

//Creates new identity element
func NewIdentity() *Identity {
	n := &Identity{}
	n.init()
	return n
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Identity) EqualLiteral(other Element) bool {
	_, ok := other.(*Identity)
	return ok
}

//Makes literal clone of element (although same() will return false)
func (c *Identity) CloneLiteral() Element {
	return NewIdentity()
}

//Turns $e$ to $a\cdot a^{-1}$ or to $a^{-1}\cdot a$ dependinf on "left". This is a step of proof.
func (c *Identity) Unannihilate(el Element, left bool) *Composite {
	if left {
		n := Compose(Inverse(el), el)
		n.setToken(c.token())
		return n
	} else {
		n := Compose(el, Inverse(el))
		n.setToken(c.token())
		return n
	}
}

package cheat1

// SHOULD NOT COMPILE
// Cheat proposed by thedeemon
// http://sober-space.livejournal.com/110937.html?thread=505689#t505689

import (
	"github.com/algebraic-brain/group_theory/gt"
	"sync"
	"testing"
)

func Test1(t *testing.T) {
	a, b := NewNamed("a"), NewNamed("b")

	if !a.same(b) {
		t.Fatal("Cheat literals should be same")
	}

	//Try to cheat verifier:
	proofForth := func(x gt.Element) gt.Element {
		return b
	}

	v := gt.VerifyForth(a, b, proofForth)

	if v {
		t.Fatal("proofForth is verified although it is not a proof")
	}
}

var lastToken = 0
var mut = &sync.Mutex{}

func tok() int {
	mut.Lock()
	token := lastToken
	lastToken++
	mut.Unlock()
	return token
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

func (el *element) ToComposite() *gt.Composite { panic("It's not Composite") }
func (el *element) ToInversed() *gt.Inversed   { panic("It's not Inversed") }
func (el *element) ToNamed() *gt.Named         { panic("It's not Named") }
func (el *element) ToIdentity() *gt.Identity   { panic("It's not Identity") }

//Checks whether one of two elements was made from other during the proof
func (el *element) Same(other gt.Element) bool {
	return el.same(other)
}

func (el *element) same(other gt.Element) bool {
	return true
}

//Composite element of group:
type Composite struct {
	element
	left  gt.Element
	right gt.Element
}

func (c *Composite) ToComposite() *gt.Composite { return nil }

//Composes elements of group
func Compose(a, b gt.Element) *Composite {
	c := &Composite{
		left:  a.CloneLiteral(),
		right: b.CloneLiteral(),
	}
	c.init()
	return c
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Composite) EqualLiteral(other gt.Element) bool {
	if o, ok := other.(*Composite); ok {
		return o.left.EqualLiteral(c.left) && o.right.EqualLiteral(c.right)
	}
	return false
}

//Makes literal clone of element (although same() will return false)
func (c *Composite) CloneLiteral() gt.Element {
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
func (c *Composite) Simplify() gt.Element {
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
func Unsimplify(el gt.Element, left bool) *Composite {
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
func (c *Composite) Map(left func(gt.Element) gt.Element, right func(gt.Element) gt.Element) *Composite {
	l := left(c.left)
	r := right(c.right)
	n := Compose(l, r)
	if l.same(c.left) && r.same(c.right) {
		n.setToken(c.token())
	}
	return n
}

//returns left element of composite
func (c *Composite) Left() gt.Element {
	return c.left.CloneLiteral()
}

//returns right element of composite
func (c *Composite) Right() gt.Element {
	return c.right.CloneLiteral()
}

//Inversed element of group:
type Inversed struct {
	element
	operand gt.Element
}

func (c *Inversed) ToInversed() *Inversed { return nil }

//Inverses element of group:
func Inverse(el gt.Element) *Inversed {
	n := &Inversed{
		operand: el,
	}
	n.init()
	return n
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Inversed) EqualLiteral(other gt.Element) bool {
	if o, ok := other.(*Inversed); ok {
		return o.operand.EqualLiteral(c.operand)
	}
	return false
}

//Makes literal clone of element (although same() will return false)
func (c *Inversed) CloneLiteral() gt.Element {
	return Inverse(c.operand.CloneLiteral())
}

//maps proofs to operand of inversion. This is a step of proof iff "f" is a step.
func (c *Inversed) Map(f func(gt.Element) gt.Element) *Inversed {
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

func (c *Named) ToNamed() *gt.Named { return nil }

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
func (c *Named) EqualLiteral(other gt.Element) bool {
	if o, ok := other.(*Named); ok {
		return o.name == c.name
	}
	return false
}

//Makes literal clone of element (although same() will return false)
func (c *Named) CloneLiteral() gt.Element {
	return NewNamed(c.name)
}

//Identity element
type Identity struct {
	element
}

func (c *Identity) ToIdentity() *Identity { return nil }

//Creates new identity element
func NewIdentity() *Identity {
	n := &Identity{}
	n.init()
	return n
}

//Checks whether two elements are equal literally (although same() may return false)
func (c *Identity) EqualLiteral(other gt.Element) bool {
	_, ok := other.(*Identity)
	return ok
}

//Makes literal clone of element (although same() will return false)
func (c *Identity) CloneLiteral() gt.Element {
	return NewIdentity()
}

//Turns $e$ to $a\cdot a^{-1}$ or to $a^{-1}\cdot a$ dependinf on "left". This is a step of proof.
func (c *Identity) Unannihilate(el gt.Element, left bool) *Composite {
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

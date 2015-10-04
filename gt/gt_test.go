package gt

import (
	"testing"
)

func TestComposeIsNotStep(t *testing.T) {
	a, b := NewNamed("a"), NewNamed("b")
	c := Compose(a, b)
	if a.Same(c) || b.Same(c) {
		t.Fatal("'Compose' is a step")
	}
}

func TestCompositeCloneLiteralIsNotStep(t *testing.T) {
	a, b := NewNamed("a"), NewNamed("b")
	c := Compose(a, b)
	d := c.CloneLiteral()
	if c.Same(d) {
		t.Fatal("'Composite.CloneLiteral' is a step")
	}
}

func TestAssociator(t *testing.T) {
	a, b, c := NewNamed("a"), NewNamed("b"), NewNamed("c")
	d := Compose(a, Compose(b, c))
	e := d.Associate()
	f := e.Unassociate()

	if !e.Same(d) || !f.Same(e) {
		t.Fatal("Associator is not a step")
	}
}

func TestAnnihilatorIsStep(t *testing.T) {
	a := NewNamed("a")
	b1 := Compose(a, Inverse(a))
	b2 := b1.Annihilate()
	b3 := b2.Unannihilate(a, false)

	if !b2.Same(b1) || !b3.Same(b2) {
		t.Fatal("Annihilator or Unannihilator is not a step (1)")
	}

	if !b3.EqualLiteral(b1) {
		t.Fatal("Annihilator and Unannihilator are not inverse of each other (1)")
	}

	c1 := Compose(Inverse(a), a)
	c2 := c1.Annihilate()
	c3 := c2.Unannihilate(a, true)

	if !c2.Same(c1) || !c3.Same(c2) {
		t.Fatal("Annihilator or Unannihilator is not a step (2)")
	}

	if !c3.EqualLiteral(c1) {
		t.Fatal("Annihilator and Unannihilator are not inverse of each other (2)")
	}

	e := NewIdentity()
	d1 := e.Unannihilate(a, true)
	d2 := d1.Annihilate()

	if !d2.EqualLiteral(e) {
		t.Fatal("Annihilator and Unannihilator are not inverse of each other (3)")
	}

	f1 := e.Unannihilate(a, false)
	f2 := f1.Annihilate()

	if !f2.EqualLiteral(e) {
		t.Fatal("Annihilator and Unannihilator are not inverse of each other (4)")
	}
}

func TestSimplifier(t *testing.T) {
	a := NewNamed("a")

	b1 := Compose(a, NewIdentity())
	b2 := b1.Simplify()
	b3 := Unsimplify(b2, false)

	if !b2.Same(b1) || !b3.Same(b2) {
		t.Fatal("Simplifier or Unsimplifier is not a step (1)")
	}

	if !b1.EqualLiteral(b3) {
		t.Fatal("Simplifier and Unsimplifier are not inverse of each other (1)")
	}

	c1 := Compose(NewIdentity(), a)
	c2 := c1.Simplify()
	c3 := Unsimplify(c2, true)

	if !c2.Same(c1) || !c3.Same(c2) {
		t.Fatal("Simplifier or Unsimplifier is not a step (2)")
	}

	if !c1.EqualLiteral(c3) {
		t.Fatal("Simplifier and Unsimplifier are not inverse of each other (2)")
	}

	d1 := Unsimplify(a, true)
	d2 := d1.Simplify()

	if !d2.EqualLiteral(a) {
		t.Fatal("Simplifier and Unsimplifier are not inverse of each other (3)")
	}

	f1 := Unsimplify(a, false)
	f2 := f1.Simplify()

	if !f2.EqualLiteral(a) {
		t.Fatal("Simplifier and Unsimplifier are not inverse of each other (4)")
	}
}

func TestCompositeMap(t *testing.T) {
	a, b := NewNamed("a"), NewNamed("b")
	c := Compose(a, b)
	f := func(el Element) Element { return el }
	d := c.Map(f, f)

	if !d.Same(c) {
		t.Fatal("'Composite.Map' is not a step when must be")
	}

	f = func(el Element) Element {
		return NewNamed("g")
	}

	h := c.Map(f, f)

	if h.Same(c) {
		t.Fatal("'Composite.Map' is a step when must not be")
	}
}

func TestLeftRightAreNotSteps(t *testing.T) {
	a, b := NewNamed("a"), NewNamed("b")
	c := Compose(a, b)

	l := c.Left()
	r := c.Right()

	if l.Same(c) || l.Same(a) || l.Same(b) || r.Same(c) || r.Same(a) || r.Same(b) {
		t.Fatal("'Left' or 'Right' are steps")
	}
}

func TestInverseIsNotStep(t *testing.T) {
	a := NewNamed("a")
	b := Inverse(a)

	if b.Same(a) {
		t.Fatal("'Inverse' is a step")
	}
}

func TestInversedCloneLiteralIsNotStep(t *testing.T) {
	a := NewNamed("a")
	c := Inverse(a)
	d := c.CloneLiteral()
	if c.Same(d) {
		t.Fatal("'Inversed.CloneLiteral' is a step")
	}
}

func TestInversedMap(t *testing.T) {
	a := NewNamed("a")
	c := Inverse(a)
	f := func(el Element) Element { return el }
	d := c.Map(f)

	if !d.Same(c) {
		t.Fatal("'Inversed.Map' is not a step when must be")
	}

	f = func(el Element) Element {
		return NewNamed("g")
	}

	h := c.Map(f)

	if h.Same(c) {
		t.Fatal("'Inversed.Map' is a step when must not be")
	}
}

func TestNamedCloneLiteralIsNotStep(t *testing.T) {
	a := NewNamed("a")
	c := a.CloneLiteral()
	if c.Same(a) {
		t.Fatal("'Named.CloneLiteral' is a step")
	}
}

func TestIdentityCloneLiteralIsNotStep(t *testing.T) {
	a := NewIdentity()
	c := a.CloneLiteral()
	if c.Same(a) {
		t.Fatal("'Identity.CloneLiteral' is a step")
	}
}

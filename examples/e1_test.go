package examples

import (
	"github.com/algebraic-brain/group_theory/gt"
	"testing"
)

func Test1(t *testing.T) {
	//Test $a \cdot b \cdot b^{-1} \cdot c = a\cdot c$
	a, b, c := gt.NewNamed("a"), gt.NewNamed("b"), gt.NewNamed("c")
	d := gt.Compose(a, gt.Compose(b, gt.Compose(gt.Inverse(b), c)))
	h := gt.Compose(a, c)

	id := func(a gt.Element) gt.Element { return a }

	proofForth := func(x gt.Element) gt.Element {
		return x.ToComposite().Map(id, func(el gt.Element) gt.Element {
			y := el.ToComposite().Associate()
			return y.Map(func(el gt.Element) gt.Element {
				return el.ToComposite().Annihilate()
			}, id).Simplify()
		})
	}

	v := gt.VerifyForth(d, h, proofForth)

	if !v {
		t.Fatal("proofForth is not verified")
	}
}

func Test2(t *testing.T) {
	a, b := gt.NewNamed("a"), gt.NewNamed("b")
	d := gt.Compose(a, b)

	//Try to cheat verifier:
	proofForth := func(x gt.Element) gt.Element {
		return x.ToComposite().Left()
	}

	v := gt.VerifyForth(d, a, proofForth)

	if v {
		t.Fatal("proofForth is verified although it is not a proof")
	}
}

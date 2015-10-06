package YourPackageName //<-- PLACE NAME OF YOUR PACKAGE HERE

import (
	"github.com/algebraic-brain/group_theory/gt"
	"testing"

	// IMPORT HERE ANYTHING YOU WANT
)

/*

YOUR CODE CAN BE PLACED HERE

*/

func TestMyCheat(t *testing.T) {
	a, b := // YOUR CODE HERE


	//Try to cheat verifier:
	proofForth := // YOUR CODE HERE

	v := gt.VerifyForth(a, b, proofForth)

	if !v {
		t.Fatal("My cheat does not work")
	}
}

/*

YOUR CODE CAN BE PLACED HERE

*/

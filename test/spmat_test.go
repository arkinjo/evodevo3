package multicell_test

import (
	"fmt"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

func TestSpMatNew(t *testing.T) {
	m := multicell.NewSpMat(200, 200, 20)
	m.Randomize(0.02)
	fmt.Println(m)
}

func TestSpMatMutate(t *testing.T) {
	m0 := multicell.NewSpMat(200, 200, 20)
	m0.Randomize(0.02)
	m1 := m0.Clone()
	if !m0.Equal(m1) {
		t.Errorf("Clone failed.")
	}

	for range 1000 {
		m0.Mutate(0.002, 0.02)
	}

	if m0.Equal(m1) {
		t.Errorf("Mutation failed.")
	}

	fmt.Printf("Densities %f %f\n", m0.Density(), m1.Density())
}

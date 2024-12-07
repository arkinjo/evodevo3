package multicell

import (
	"math"
)

type Cell struct {
	States [][]float64
	Pave   []float64
	Pvar   []float64
}

func (c *Cell) Left() Vec {
	return c.Pave[0 : len(c.Pave)/4]
}

func (c *Cell) Top() Vec {
	return c.Pave[len(c.Pave)/4 : len(c.Pave)/2]
}

func (c *Cell) Right() Vec {
	return c.Pave[len(c.Pave)/2 : len(c.Pave)*3/4]
}

func (c *Cell) Bottom() Vec {
	return c.Pave[len(c.Pave)*3/4:]
}

func (s *Setting) NewCell() Cell {
	states := make([][]float64, s.Num_layers)
	for i, nc := range s.Num_components {
		states[i] = NewVec(nc, 1.0)
	}
	pave := make([]float64, s.Num_components[s.Num_layers-1])
	pvar := make([]float64, s.Num_components[s.Num_layers-1])

	return Cell{
		States: states,
		Pave:   pave,
		Pvar:   pvar,
	}
}

func (c *Cell) Initialize() {
	for l, state := range c.States {
		for i := range state {
			c.States[l][i] = 1.0
		}
	}
}

func (c *Cell) Dev_step(s *Setting, g Genome, istep int, env Vec) float64 {
	for l := 1; l < s.Num_layers; l++ {
		va := make(Vec, s.Num_components[l])
		vt := make(Vec, s.Num_components[l])
		for k, mat := range g[l] {
			vu := make(Vec, s.Num_components[k])
			if l == 1 && k == 0 {
				DiffVecs(c.States[k], c.Pave, vu)
			} else {
				vu = c.States[k]
			}
			MultSpMatVec(mat, vu, vt)
			AddVecs(va, vt, va)
		}
		if l < s.Num_layers-1 {
			ApplyFVec(LCatan, s.Omega[l], va, c.States[l])
		} else {
			ApplyFVec(math.Tanh, s.Omega[l], va, c.States[l])
		}
	}
	if istep == 0 {
		for i, p := range c.States[s.Num_layers-1] {
			c.Pave[i] = p
			c.Pvar[i] = 1.0
		}
	} else { // exponential moving average/variance
		for i, p := range c.States[s.Num_layers-1] {
			d := p - c.Pave[i]
			incr := s.Alpha * d
			c.Pave[i] += incr
			c.Pvar[i] = (1 - s.Alpha) * (c.Pvar[i] + d*incr)
		}
	}
	dev := 0.0
	for _, d := range c.Pvar {
		dev += d
	}
	return dev / float64(s.Num_components[s.Num_layers-1])
}

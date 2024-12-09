package multicell

import (
	//	"fmt"
	"math"
)

type Cell struct {
	States [][]float64
	Pave   []float64
	Pvar   []float64

	El Vec // -> cell[i-1][j].Pr || Envs.Lefts[j]
	Et Vec // -> cell[i][j+1].Pb || Envs.Tops[i]
	Er Vec // -> cell[i+1][j].Pl || Envs.Rights[j]
	Eb Vec // -> cell[i][j-1].Pt || Envs.Bottoms[i]

	Pl, Pt, Pr, Pb Vec // -> parts of Pave
}

func (s *Setting) NewCell() Cell {
	states := make([][]float64, s.Num_layers)
	for i, nc := range s.Num_components {
		states[i] = NewVec(nc, 1.0)
	}

	lenP := s.Num_components[s.Num_layers-1]
	pave := make([]float64, lenP)
	pvar := make([]float64, lenP)

	cell := Cell{
		States: states,
		Pave:   pave,
		Pvar:   pvar,
	}
	cell.Initialize(s)
	return cell
}

func (c *Cell) Initialize(s *Setting) {
	for l := range c.States {
		SetVec(c.States[l], 1.0)
	}
	SetVec(c.Pave, 1.0)
	SetVec(c.Pvar, 1.0)
	lenP4 := s.Num_components[s.Num_layers-1] / 4
	c.Pl = c.Pave[:lenP4]
	c.Pt = c.Pave[lenP4 : lenP4*2]
	c.Pr = c.Pave[lenP4*2 : lenP4*3]
	c.Pb = c.Pave[lenP4*3:]
}

func (c *Cell) Set_env() {
	k := 0
	for _, v := range c.El {
		c.States[0][k] = v
		k++
	}
	for _, v := range c.Et {
		c.States[0][k] = v
		k++
	}
	for _, v := range c.Er {
		c.States[0][k] = v
		k++
	}
	for _, v := range c.Eb {
		c.States[0][k] = v
		k++
	}
}

func (c *Cell) Dev_step(s *Setting, g Genome, istep int) float64 {
	c.Set_env()
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

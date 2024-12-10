package multicell

import (
	//	"fmt"
	"math"
)

type Cell struct {
	E    []Vec       // points to neighboring cell face or environment
	M    [][]float64 // middle layers
	P    []Vec       // output layer
	Pave []Vec
	Pvar []Vec
}

func (s *Setting) NewCell() Cell {
	e := make([]Vec, NumFaces) //
	p := make([]Vec, NumFaces)
	pave := make([]Vec, NumFaces)
	pvar := make([]Vec, NumFaces)

	for i := range NumFaces {
		p[i] = NewVec(s.LenFace, 1.0)
		pave[i] = NewVec(s.LenFace, 1.0)
		pvar[i] = NewVec(s.LenFace, 1.0)
	}

	m := make([][]float64, s.NumLayers)
	for i, nc := range s.LenLayer {
		m[i] = NewVec(nc, 1.0)
	}

	return Cell{
		E:    e,
		M:    m,
		P:    p,
		Pave: pave,
		Pvar: pvar}
}

func (c *Cell) Initialize(s *Setting) {
	for i := range NumFaces {
		SetVec(c.E[i], 1.0)
		SetVec(c.P[i], 1.0)
		SetVec(c.Pave[i], 1.0)
		SetVec(c.Pvar[i], 1.0)
	}
	for l := range s.NumLayers {
		SetVec(c.M[l], 1.0)
	}
}

func (c *Cell) DevStep(s *Setting, g Genome, istep int) float64 {
	v0 := make(Vec, s.LenFace)
	v1 := make(Vec, s.LenLayer[0])
	s0 := make(Vec, s.LenLayer[0])
	for i, vi := range c.E {
		DiffVecs(v0, vi, c.Pave[i])
		MultSpMatVec(v1, g.E[i], v0)
		AddVecs(s0, s0, v1)
	}

	for l := 0; l < s.NumLayers; l++ {
		va := make(Vec, s.LenLayer[l])
		vt := make(Vec, s.LenLayer[l])
		if l == 0 {
			AddVecs(va, va, s0)
		}
		for k, mat := range g.M[l] {
			MultSpMatVec(vt, mat, c.M[k])
			AddVecs(va, va, vt)
		}
		ApplyFVec(c.M[l], LCatan, s.Omega[l], va)
	}

	w0 := make(Vec, s.LenFace)
	for i := range c.P {
		MultSpMatVec(w0, g.P[i], c.M[s.NumLayers-1])
		ApplyFVec(c.P[i], math.Tanh, s.OmegaP, w0)
	}
	if istep == 0 {
		for i, p := range c.P {
			copy(c.Pave[i], p)
			SetVec(c.Pvar[i], 1.0)
		}
	} else { // exponential moving average/variance
		for i, p := range c.P {
			for j, v := range p {
				d := v - c.Pave[i][j]
				incr := s.Alpha * d
				c.Pave[i][j] += incr
				c.Pvar[i][j] = (1 - s.Alpha) * (c.Pvar[i][j] + d*incr)
			}
		}
	}
	dev := 0.0
	for _, p := range c.Pvar {
		for _, d := range p {
			dev += d
		}
	}
	return dev / float64(s.LenFace*NumFaces)
}

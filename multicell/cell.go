package multicell

import (
	//	"fmt"
	"log"
)

type Cell struct {
	Id     int           // Identifier within an individual
	Facing [NumFaces]int // Facing Cell's Id; -1 if none.
	E      []Vec         // points to neighboring cell face or environment
	S      [][]float64   // middle and output layers
	Pave   Vec
	Pvar   Vec
}

func (s *Setting) NewCell(id int) Cell {
	var facing [NumFaces]int
	for i := range NumFaces {
		facing[i] = -1
	}
	e := make([]Vec, NumFaces) //
	pave := NewVec(s.LenLayer[s.NumLayers-1], 1.0)
	pvar := NewVec(s.LenLayer[s.NumLayers-1], 1.0)

	m := make([][]float64, s.NumLayers)
	for i, nc := range s.LenLayer {
		m[i] = NewVec(nc, 1.0)
	}

	return Cell{
		Id:     id,
		Facing: facing,
		E:      e,
		S:      m,
		Pave:   pave,
		Pvar:   pvar}
}

func (c *Cell) Initialize(s *Setting) {
	for l := range s.NumLayers {
		SetVec(c.S[l], 1.0)
	}
}

func (c *Cell) Left(s *Setting) Vec {
	return c.Pave[:s.LenFace]
}

func (c *Cell) Top(s *Setting) Vec {
	return c.Pave[s.LenFace : s.LenFace*2]
}

func (c *Cell) Right(s *Setting) Vec {
	return c.Pave[s.LenFace*2 : s.LenFace*3]
}

func (c *Cell) Bottom(s *Setting) Vec {
	return c.Pave[s.LenFace*3:]
}

func (c *Cell) Face(s *Setting, iface int) Vec {
	var v Vec
	switch iface {
	case Left:
		v = c.Left(s)
	case Top:
		v = c.Top(s)
	case Right:
		v = c.Right(s)
	case Bottom:
		v = c.Bottom(s)
	default:
		log.Fatal("(*cell).Face: Unknown face")
	}
	return v
}

func (c *Cell) OppositeFace(s *Setting, iface int) Vec {
	var v Vec
	switch iface {
	case Left:
		v = c.Right(s)
	case Top:
		v = c.Bottom(s)
	case Right:
		v = c.Left(s)
	case Bottom:
		v = c.Top(s)
	default:
		log.Fatal("(*cell).OppositeFace: Unknown face")
	}
	return v
}

func (c *Cell) DevStep(s *Setting, g Genome, istep int) float64 {
	v0 := make(Vec, s.LenFace)
	v1 := make(Vec, s.LenLayer[0])
	s0 := make(Vec, s.LenLayer[0])
	for i, vi := range c.E {
		DiffVecs(v0, vi, c.Face(s, i))
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
			MultSpMatVec(vt, mat, c.S[k])
			AddVecs(va, va, vt)
		}
		if l < s.NumLayers-1 {
			ApplyFVec(c.S[l], LCatan(s.Omega[l]), va)
		} else {
			ApplyFVec(c.S[l], Tanh(s.Omega[l]), va)
		}
	}

	if istep == 0 {
		copy(c.Pave, c.S[s.NumLayers-1])
		SetVec(c.Pvar, 1.0)
	} else { // exponential moving average/variance
		for i, v := range c.S[s.NumLayers-1] {
			d := v - c.Pave[i]
			incr := s.Alpha * d
			c.Pave[i] += incr
			c.Pvar[i] = (1 - s.Alpha) * (c.Pvar[i] + d*incr)
		}
	}
	dev := 0.0
	for _, d := range c.Pvar {
		dev += d
	}
	return dev / float64(s.LenLayer[s.NumLayers-1])
}

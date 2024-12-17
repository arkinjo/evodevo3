package multicell

import (
	//	"fmt"
	"log"
)

type Cell struct {
	Id     int           // Identifier within an individual
	Facing [NumFaces]int // Facing Cell's Id; -1 if none.
	E      []Vec         // points to neighboring cell face or environment
	S      []Vec         // middle and output layers
	P      Vec
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

	m := make([]Vec, s.NumLayers)
	for i, nc := range s.LenLayer {
		m[i] = NewVec(nc, 1.0)
	}

	return Cell{
		Id:     id,
		Facing: facing,
		E:      e,
		S:      m,
		P:      pave,
		Pvar:   pvar}
}

func (c *Cell) Initialize(s *Setting) {
	for l := range s.NumLayers {
		c.S[l].SetAll(1.0)
	}
}

func (c *Cell) Left(s *Setting) Vec {
	return c.P[:s.LenFace]
}

func (c *Cell) Top(s *Setting) Vec {
	return c.P[s.LenFace : s.LenFace*2]
}

func (c *Cell) Right(s *Setting) Vec {
	return c.P[s.LenFace*2 : s.LenFace*3]
}

func (c *Cell) Bottom(s *Setting) Vec {
	return c.P[s.LenFace*3:]
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
	for i, ei := range c.E {
		v0.Diff(ei, c.Face(s, i))
		g.E[i].MultVec(v0, v1)
		s0.Acc(v1)
	}

	for l, tl := range s.Topology {
		va := make(Vec, s.LenLayer[l])
		vt := make(Vec, s.LenLayer[l])
		if l == 0 {
			va.Acc(s0)
		}
		for k := range tl {
			mat := g.M[IntPair{l, k}]
			mat.MultVec(c.S[k], vt)
			va.Acc(vt)
		}
		var afunc func(float64) float64
		if l < s.NumLayers-1 {
			afunc = LCatan(g.W[l])
		} else {
			afunc = Tanh(g.W[l])
		}
		c.S[l].ApplyFVec(afunc, va)
	}

	if istep == 0 {
		copy(c.P, c.S[s.NumLayers-1])
		c.Pvar.SetAll(1.0)
	} else { // exponential moving average/variance
		for i, v := range c.S[s.NumLayers-1] {
			d := v - c.P[i]
			incr := s.Alpha * d
			c.P[i] += incr
			c.Pvar[i] = (1 - s.Alpha) * (c.Pvar[i] + d*incr)
		}
	}
	dev := 0.0
	for _, d := range c.Pvar {
		dev += d
	}
	return dev / float64(s.LenLayer[s.NumLayers-1])
}

package multicell

import (
	//	"fmt"
	"log"
	"slices"
)

type Cell struct {
	Id     int           // Identifier within an individual
	Facing [NumFaces]int // Facing Cell's Id; -1 if none.
	Cue    []Vec         // points to neighboring cell face or environment
	S      []Vec         // state vectors
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
	pvar := NewVec(s.LenLayer[s.NumLayers-1], 0.0)

	m := make([]Vec, s.NumLayers)
	for i, nc := range s.LenLayer {
		m[i] = NewVec(nc, 1.0).AddNoise(s.NumBlocks*s.LenBlock, 1)
	}

	return Cell{
		Id:     id,
		Facing: facing,
		Cue:    e,
		S:      m,
		Pave:   pave,
		Pvar:   pvar}
}

func (c *Cell) Initialize(s *Setting) {
	for l := range s.NumLayers {
		c.S[l].SetAll(1.0)
	}
	c.Pave.SetAll(1.0)
	c.Pvar.SetAll(0.0)
}

// all internal states into one vector.
func (c *Cell) ToVec() Vec {
	return slices.Concat(c.S[1 : len(c.S)-1]...)
}

func (c *Cell) Left(s *Setting) Vec {
	return c.S[s.NumLayers-1][:s.LenFace]
}

func (c *Cell) Top(s *Setting) Vec {
	return c.S[s.NumLayers-1][s.LenFace : s.LenFace*2]
}

func (c *Cell) Right(s *Setting) Vec {
	return c.S[s.NumLayers-1][s.LenFace*2 : s.LenFace*3]
}

func (c *Cell) Bottom(s *Setting) Vec {
	return c.S[s.NumLayers-1][s.LenFace*3:]
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
	s.Topology.EachRow(func(l int, tl map[int]float64) {
		va := make(Vec, s.LenLayer[l])
		if l == 0 {
			s0 := slices.Concat(c.Cue...)
			va.Diff(s0, c.S[s.NumLayers-1])
		}
		for k := range tl {
			va.MultSpMatVec(g.M[l][k], c.S[k]) // va is accumulated.
		}
		if with_bias {
			va.Acc(g.B[l])
		}
		afunc := LCatan(s.Omega[l])
		if l == s.NumLayers-1 {
			afunc = CStep1(s.Omega[l])
		}
		c.S[l].ApplyFVec(afunc, va)
	})

	for i, v := range c.S[s.NumLayers-1] {
		d := v - c.Pave[i]
		incr := s.Alpha * d
		c.Pave[i] += incr
		c.Pvar[i] = (1 - s.Alpha) * (c.Pvar[i] + d*incr)
	}

	dev := 0.0
	for _, d := range c.Pvar {
		dev += d
	}
	return dev / float64(s.LenLayer[s.NumLayers-1])
}

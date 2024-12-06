package multicell

type Cell struct {
	States [][]float64
	Pave   []float64
	Pvar   []float64
}

func (s *Setting) Left(c *Cell) Vec {
	return c.States[s.Num_layers-1][0 : s.Num_components[0]/4]
}

func (s *Setting) Top(c *Cell) Vec {
	return c.States[s.Num_layers-1][s.Num_components[0]/4 : s.Num_components[0]/2]
}

func (s *Setting) Right(c *Cell) Vec {
	return c.States[s.Num_layers-1][s.Num_components[0]/2 : s.Num_components[0]*3/4]
}

func (s *Setting) Bottom(c *Cell) Vec {
	return c.States[s.Num_layers-1][s.Num_components[0]*3/4:]
}

func (s *Setting) NewCell() Cell {
	states := make([][]float64, s.Num_layers)
	for i, nc := range s.Num_components {
		states[i] = make([]float64, nc)
		for j, _ := range states[i] {
			states[i][j] = 1.0
		}
	}
	pave := make([]float64, s.Num_components[s.Num_layers-1])
	pvar := make([]float64, s.Num_components[s.Num_layers-1])

	return Cell{
		States: states,
		Pave:   pave,
		Pvar:   pvar,
	}
}

func (c *Cell) Dev_step(s *Setting, g Genome, istep int, env Vec) float64 {
	for l := 1; l < s.Num_layers; l++ {
		va := make(Vec, s.Num_components[l])
		vt := make(Vec, s.Num_components[l])
		for k, mat := range g[l] {
			MultSpMatVec(mat, c.States[k], vt)
			AddVecs(va, vt, va)
		}
		ApplyFVec(s.Afuncs[l], s.Omega[l], va, c.States[l])
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

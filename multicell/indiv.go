package multicell

import (
	//	"fmt"
	"math"
	"slices"
)

type Individual struct {
	Id       int
	Mom_id   int
	Dad_id   int
	Genome   Genome
	Cells    [][]Cell
	Envs     Cell_envs
	Ndev     int
	Mismatch float64
	Fitness  float64
}

func (s *Setting) NewIndividual(id int, env Environment) Individual {
	cells := make([][]Cell, s.Num_cell_x)
	envs := s.NewCell_envs(env)
	for i := range s.Num_cell_x {
		cells[i] = make([]Cell, s.Num_cell_y)
		for j := range s.Num_cell_y {
			cells[i][j] = s.NewCell()
		}
	}

	return Individual{
		Id:       id,
		Mom_id:   -1,
		Dad_id:   -1,
		Genome:   s.NewGenome(),
		Cells:    cells,
		Envs:     envs,
		Ndev:     0,
		Mismatch: math.Inf(0),
		Fitness:  0,
	}
}

func (indiv *Individual) Selected_pheno(s *Setting) Vec {
	var p Vec
	for j := range s.Num_cell_y {
		p = append(p, indiv.Cells[0][j].Left()...)
	}

	return p
}

func (indiv *Individual) Get_cell_env(s *Setting, i, j int) Vec {
	var left, right, top, bottom Vec

	for j := range s.Num_cell_y {
		if i == 0 {
			left = append(left, indiv.Envs.Rights[j]...)
		} else {
			left = append(left, indiv.Cells[i-1][j].Right()...)
		}
		if i == s.Num_cell_x-1 {
			right = append(right, indiv.Envs.Lefts[j]...)
		} else {
			right = append(right, indiv.Cells[i+1][j].Left()...)
		}
	}

	for i := range s.Num_cell_x {
		if j == 0 {
			bottom = append(bottom, indiv.Envs.Tops[i]...)
		} else {
			bottom = append(bottom, indiv.Cells[i][j].Top()...)
		}
		if j == s.Num_cell_y-1 {
			top = append(top, indiv.Envs.Bottoms[i]...)
		} else {
			top = append(top, indiv.Cells[i][j].Bottom()...)
		}
	}
	return slices.Concat(left, top, right, bottom)
}

func (indiv *Individual) Get_mismatch(s *Setting, env Environment) float64 {
	selenv := s.Selecting_env(env)
	selphen := indiv.Selected_pheno(s)
	dev := 0.0
	for i, e := range selenv {
		dev += math.Abs(e - selphen[i])
	}

	return dev / float64(len(selenv))
}

func (indiv *Individual) Develop(s *Setting, env Environment) Individual {
	istep := 0
	for istep = range s.Num_dev {
		dev := 0.0
		for i := range s.Num_cell_x {
			for j := range s.Num_cell_y {
				env := indiv.Get_cell_env(s, i, j)
				dev += indiv.Cells[i][j].Dev_step(s, indiv.Genome, istep, env)
			}
		}
		if dev < s.Conv_dev {
			break
		}
	}

	indiv.Ndev = istep
	indiv.Mismatch = indiv.Get_mismatch(s, env)

	if istep < s.Num_dev {
		indiv.Fitness = math.Exp(-indiv.Mismatch * s.Selstrength)
	} else {
		indiv.Fitness = 0
	}
	return *indiv
}

func (s *Setting) MateIndividuals(indiv0, indiv1 Individual, env Environment) (Individual, Individual) {
	geno0, geno1 := s.MateGenomes(indiv0.Genome, indiv1.Genome)
	kid0 := s.NewIndividual(-1, env)
	kid1 := s.NewIndividual(-2, env)

	kid0.Genome = geno0
	kid0.Mom_id = indiv0.Id
	kid0.Dad_id = indiv1.Id

	kid1.Genome = geno1
	kid0.Mom_id = indiv1.Id
	kid0.Dad_id = indiv0.Id

	return kid0, kid1
}

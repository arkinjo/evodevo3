package main

import (
	"fmt"
	"github.com/arkinjo/evodevo3/multicell"
	//	"math"
)

type Hoge interface {
	Gaga() float64
}

type Cl1 struct {
	a float64
}

func (c Cl1) Gaga() float64 {
	return c.a
}

type Cl2 struct {
	b string
}

func (c Cl2) Gaga() float64 {
	return float64(len(c.b))
}

func main() {
	s := multicell.Get_default_setting("hoge", 5, 15)
	s.Set_omega()
	for i, tl := range s.Topology {
		for j, d := range tl {
			fmt.Println("Topology ", i, " ", j, " ", d)
		}
	}
	s.Max_pop = 500
	env := s.NewEnvironment()
	pop := s.NewPopulation(env)
	pop.Evolve(s, 200, env)
	fmt.Println(pop)
	fmt.Println(s.Selecting_env(env))
	fmt.Println(pop.Indivs[0].Selected_pheno(s))

}

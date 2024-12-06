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
	env := s.NewEnvironment()
	indiv := s.NewIndividual(1, env)
	indiv.Develop(s, env)
	fmt.Println(len(s.Selecting_env(env)))
	fmt.Println(indiv)
}

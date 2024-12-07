package multicell

import (
	"fmt"
	"math/rand/v2"
)

type Population struct {
	Igen   int // generation
	Env    Environment
	Indivs []Individual
}

func (s *Setting) NewPopulation(env Environment) Population {
	var indivs []Individual
	for id := range s.Max_pop {
		indivs = append(indivs, s.NewIndividual(id, env))
	}
	return Population{
		Igen:   0,
		Indivs: indivs}
}

func (pop *Population) Get_max_fitness() float64 {
	f := 0.0
	for _, indiv := range pop.Indivs {
		if f < indiv.Fitness {
			f = indiv.Fitness
		}
	}
	return f
}

func (pop *Population) Develop(s *Setting, igen int, env Environment) {
	pop.Igen = igen
	ch := make(chan Individual)

	for _, indiv := range pop.Indivs {
		go func(indiv Individual) {
			ch <- indiv.Develop(s, env)
		}(indiv)
	}

	for i := range pop.Indivs {
		pop.Indivs[i] = <-ch
	}
}

func (pop *Population) Select(s *Setting, env Environment) Population {
	var indivs []Individual
	maxfit := pop.Get_max_fitness()
	npop := 0
	for {
		i := rand.IntN(s.Max_pop)
		wfit := pop.Indivs[i].Fitness / maxfit

		if rand.Float64() < wfit {
			indivs = append(indivs, pop.Indivs[i])
			npop++
		}
		if npop == s.Max_pop {
			break
		}
	}
	return Population{
		Igen:   pop.Igen,
		Env:    pop.Env,
		Indivs: indivs}
}

func (pop *Population) Reproduce(s *Setting, env Environment) Population {
	var kids []Individual

	for i := range len(pop.Indivs) {
		if i == 0 {
			continue
		}
		kid0, kid1 := s.MateIndividuals(pop.Indivs[i-1], pop.Indivs[i], env)
		kid0.Id = i - 1
		kid1.Id = i
		kids = append(kids, kid0)
		kids = append(kids, kid1)
	}

	return Population{
		Igen:   pop.Igen + 1,
		Env:    env,
		Indivs: kids}
}

func (pop0 *Population) Evolve(s *Setting, maxgen int, env Environment) Population {
	pop := *pop0
	for igen := range maxgen {
		fmt.Println("Evolve: ", igen)
		pop.Develop(s, igen, env)
		pop = pop.Select(s, env)
		pop = pop.Reproduce(s, env)
	}
	return pop
}

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

type PopStats struct {
	Mismatch float64
	Fitness  float64
	Ndev     float64
	Nparents int
}

func (pop *Population) Get_popstats() PopStats {
	mismatch := 0.0
	fitness := 0.0
	ndev := 0.0
	npop := float64(len(pop.Indivs))
	npar := make(map[int]bool)
	for _, indiv := range pop.Indivs {
		mismatch += indiv.Mismatch
		fitness += indiv.Fitness
		ndev += float64(indiv.Ndev)
		npar[indiv.Id] = true
	}
	return PopStats{
		Mismatch: mismatch / npop,
		Fitness:  fitness / npop,
		Ndev:     ndev / npop,
		Nparents: len(npar)}
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

func (pop *Population) Develop(s *Setting, igen int, selenv Environment) {
	pop.Igen = igen
	ch := make(chan Individual)

	for _, indiv := range pop.Indivs {
		go func(indiv Individual) {
			ch <- indiv.Develop(s, selenv)
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

func (pop0 *Population) Evolve(s *Setting, env Environment, maxgen int) Population {
	pop := *pop0
	selenv := s.Selecting_env(env)
	for igen := range maxgen {
		pop.Develop(s, igen+1, selenv)
		pop = pop.Select(s, env)
		stats := pop.Get_popstats()
		fmt.Printf("%d\t%e\t%e\t%e\t%d\n",
			pop.Igen, stats.Mismatch, stats.Fitness, stats.Ndev, stats.Nparents)
		pop = pop.Reproduce(s, env)
	}
	return pop
}

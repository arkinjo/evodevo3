package multicell

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"os"
)

type Population struct {
	Iepoch int // epoch
	Igen   int // generation
	Indivs []Individual
}

type PopStats struct {
	Mismatch float64
	Fitness  float64
	Ndev     float64
	Nparents int
}

func (pop *Population) GetPopStats() PopStats {
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
	for id := range s.MaxPopulation {
		indiv := s.NewIndividual(id, env)
		indiv.Genome = s.NewGenome()
		indivs = append(indivs, indiv)
	}
	return Population{
		Iepoch: 0,
		Igen:   0,
		Indivs: indivs}
}

func (pop *Population) GetMaxFitness() float64 {
	f := 0.0
	for _, indiv := range pop.Indivs {
		if f < indiv.Fitness {
			f = indiv.Fitness
		}
	}
	return f
}

func (pop *Population) Develop(s *Setting, selenv Environment) {
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

func (pop *Population) Select(s *Setting) Population {
	var indivs []Individual
	maxfit := pop.GetMaxFitness()
	npop := 0
	for {
		i := rand.IntN(s.MaxPopulation)
		wfit := pop.Indivs[i].Fitness / maxfit

		if rand.Float64() < wfit {
			indivs = append(indivs, pop.Indivs[i])
			npop++
		}
		if npop == s.MaxPopulation {
			break
		}
	}
	return Population{
		Iepoch: pop.Iepoch,
		Igen:   pop.Igen,
		Indivs: indivs}
}

func (pop *Population) Reproduce(s *Setting, env Environment) Population {
	var kids []Individual

	for i, indiv := range pop.Indivs {
		if i == 0 {
			continue
		}

		kid0, kid1 := s.MateIndividuals(pop.Indivs[i-1], indiv, env)
		kid0.Id = i - 1
		kid1.Id = i
		kids = append(kids, kid0)
		kids = append(kids, kid1)
	}

	return Population{
		Iepoch: pop.Iepoch,
		Igen:   pop.Igen + 1,
		Indivs: kids}
}

func (pop0 *Population) Evolve(s *Setting, maxgen int, env Environment) Population {
	pop := *pop0

	selenv := s.SelectingEnv(env)
	for igen := range maxgen {
		pop.Igen = igen
		pop.Develop(s, selenv)
		pop = pop.Select(s)
		stats := pop.GetPopStats()
		fmt.Printf("%d\t%d\t%e\t%e\t%e\t%d\n",
			pop.Iepoch, pop.Igen,
			stats.Mismatch, stats.Fitness, stats.Ndev, stats.Nparents)
		if s.ProductionRun {
			pop.Dump(s)
		}
		pop = pop.Reproduce(s, env)
	}
	pop.Igen = maxgen
	pop.Develop(s, selenv)
	pop.Dump(s)
	return pop
}

func (pop *Population) Initialize(s *Setting, env Environment) {
	for i := range pop.Indivs {
		pop.Indivs[i].Initialize(s, env)
	}
}

func (s *Setting) TrajectoryFilename(iepoch, igen int) string {
	filename := fmt.Sprintf("%s/%s_%2.2d_%3.3d.traj.gz", s.Outdir, s.Basename, iepoch, igen)
	return filename
}

func (pop *Population) Dump(s *Setting) string {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen)
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()
	foutz, err := gzip.NewWriterLevel(fout, gzip.BestSpeed)
	JustFail(err)
	defer foutz.Close()

	encoder := gob.NewEncoder(foutz)
	encoder.Encode(pop)
	return filename
}

func (s *Setting) LoadPopulation(filename string, env Environment) Population {
	pop := s.NewPopulation(env)
	fin, err := os.Open(filename)
	JustFail(err)
	defer fin.Close()

	finz, err := gzip.NewReader(fin)
	JustFail(err)
	defer finz.Close()

	decoder := gob.NewDecoder(finz)
	err = decoder.Decode(&pop)
	JustFail(err)
	return pop
}

func (pop *Population) GenomeVecs(s *Setting) []Vec {
	vecs := make([]Vec, len(pop.Indivs))
	for i, indiv := range pop.Indivs {
		vecs[i] = indiv.Genome.ToVec(s)
	}
	return vecs
}

func (pop *Population) CueVecs(s *Setting) []Vec {
	vecs := make([]Vec, len(pop.Indivs))
	for i, indiv := range pop.Indivs {
		vecs[i] = indiv.CueVec(s)
	}
	return vecs
}

func (pop *Population) PhenoVecs() []Vec {
	vecs := make([]Vec, len(pop.Indivs))
	for i, indiv := range pop.Indivs {
		vecs[i] = indiv.SelectedPhenotype()
	}
	return vecs
}

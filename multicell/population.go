package multicell

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"sort"
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

func (stats PopStats) Print(iepoch, igen int) {
	fmt.Printf("%d\t%d\t%e\t%e\t%e\t%d\n",
		iepoch, igen,
		stats.Mismatch, stats.Fitness, stats.Ndev, stats.Nparents)
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
	ch := make(chan Individual, 2)

	for i := 1; i < len(pop.Indivs); i += 2 {
		go func(mom, dad Individual) {
			kid0, kid1 := s.MateIndividuals(mom, dad, env)
			ch <- kid0
			ch <- kid1
		}(pop.Indivs[i-1], pop.Indivs[i])
	}

	var kids []Individual
	for i := 1; i < len(pop.Indivs); i += 2 {
		kid := <-ch
		kid.Id = i - 1
		kids = append(kids, kid)

		kid = <-ch
		kid.Id = i
		kids = append(kids, kid)
	}

	return Population{
		Iepoch: pop.Iepoch,
		Igen:   pop.Igen + 1,
		Indivs: kids}
}

func (pop0 *Population) Evolve(s *Setting, env Environment) Population {
	pop := *pop0

	selenv := s.SelectingEnv(env)
	for igen := range s.MaxGeneration {
		pop.Igen = igen
		pop.Develop(s, selenv)
		if s.ProductionRun { // Dump before Selection
			pop.Dump(s)
		}
		pop = pop.Select(s)
		stats := pop.GetPopStats()
		stats.Print(pop.Iepoch, pop.Igen)
		pop = pop.Reproduce(s, env)
	}
	pop.Igen = s.MaxGeneration
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
	filename := fmt.Sprintf("%s/%s_%2.2d_%3.3d.traj", s.Outdir, s.Basename, iepoch, igen)
	return filename
}

// Dump the Population in a gzipped binary file.
func (pop *Population) Dump(s *Setting) string {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen) + ".gz"
	fout, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	JustFail(err)
	defer fout.Close()
	foutz, err := gzip.NewWriterLevel(fout, gzip.BestSpeed)
	JustFail(err)
	defer foutz.Close()

	encoder := gob.NewEncoder(foutz)
	encoder.Encode(pop)
	log.Printf("Trajectory Dump saved in: %s\n", filename)
	return filename
}

func (pop *Population) DumpJSON(s *Setting) string {
	filename := s.TrajectoryFilename(pop.Iepoch, pop.Igen) + ".json"
	json, err := json.MarshalIndent(pop, "", "  ")
	JustFail(err)
	os.WriteFile(filename, json, 0644)
	log.Printf("Trajectory JSON saved in: %s\n", filename)
	return filename
}

func (s *Setting) LoadPopulation(filename string, env Environment) Population {
	log.Printf("Load population from: %s\n", filename)
	fin, err := os.Open(filename)
	JustFail(err)
	defer fin.Close()

	finz, err := gzip.NewReader(fin)
	JustFail(err)
	defer finz.Close()

	decoder := gob.NewDecoder(finz)

	var pop Population
	err = decoder.Decode(&pop)
	JustFail(err)

	sort.SliceStable(pop.Indivs, func(i, j int) bool {
		return pop.Indivs[i].Id < pop.Indivs[j].Id
	})

	s.MaxPopulation = len(pop.Indivs)
	pop.Initialize(s, env)
	return pop
}

func (s *Setting) LoadPopulationJSON(filename string, env Environment) Population {
	log.Printf("Load population JSON from: %s\n", filename)
	buffer, err := os.ReadFile(filename)
	JustFail(err)
	var pop Population
	err = json.Unmarshal(buffer, &pop)
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

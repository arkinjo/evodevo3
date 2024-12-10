package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
	Estart  int
	Eend    int
	Ngen    int
}

func GetSetting() Simulation {
	maxpopP := flag.Int("popsize", 500, "population size")
	nenvsP := flag.Int("num_envs", 50, "number of environments")
	envsfileP := flag.String("envs", "", "saved environments JSON file")
	mkenvsP := flag.String("make_envs", "", "Make brand new environments")
	modelP := flag.String("model", "Full", "Model name")
	trajDirP := flag.String("trajdir", "traj", "Directory for trajectory files")
	envStartP := flag.Int("env_start", 0, "starting environment (0, 1, ...)")
	envEndP := flag.Int("env_end", 20, "ending environment")
	ngenP := flag.Int("ngen", 200, "number of generations per epoch")
	prodP := flag.Bool("production", false, "true if production run")

	flag.Parse()

	s := multicell.GetDefaultSetting()
	s.MaxPopulation = *maxpopP
	s.Outdir = *trajDirP
	s.SetModel(*modelP)
	s.ProductionRun = *prodP

	var envs []multicell.Environment

	if *mkenvsP == "" && *envsfileP == "" {
		panic("use -envs or -make_envs")
	} else if *mkenvsP != "" {
		s.SaveEnvs(*mkenvsP, *nenvsP)
		log.Printf("Environments saved in: %s; exitting...\n", *mkenvsP)
		os.Exit(0)
	} else if *envsfileP != "" {
		envs = s.LoadEnvs(*envsfileP)
	}

	if *envStartP >= *envEndP {
		log.Printf("The starting environment must be < the ending environment (%d)", *envEndP)
		os.Exit(1)
	}
	if *envStartP >= len(envs) {
		log.Printf("The starting envirnment must be < %d", len(envs))
		os.Exit(1)
	}

	if *envEndP > len(envs) {
		log.Printf("The ending envirnment must be <= %d", len(envs))
		os.Exit(1)
	}

	return Simulation{
		Setting: s,
		Envs:    envs,
		Estart:  *envStartP,
		Eend:    *envEndP,
		Ngen:    *ngenP}

}

func main() {
	t0 := time.Now()
	sim := GetSetting()
	fmt.Println(sim.Setting, sim.Estart, sim.Eend)
	pop := sim.Setting.NewPopulation(sim.Envs[0])
	for iepoch := sim.Estart; iepoch < sim.Eend; iepoch++ {
		pop.Iepoch += 1
		pop = pop.Evolve(sim.Setting, sim.Ngen, sim.Envs[iepoch])
	}
	fmt.Println("Time: ", time.Since(t0))
}

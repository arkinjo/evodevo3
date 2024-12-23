package main

import (
	"flag"
	//	"fmt"
	"log"
	"os"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
	Pop     multicell.Population
	Estart  int
	Eend    int
}

func GetSetting() Simulation {
	maxpopP := flag.Int("popsize", 500, "population size")
	nenvsP := flag.Int("num_envs", 50, "number of environments")
	envsfileP := flag.String("envs", "", "saved environments JSON file")
	resfileP := flag.String("restart", "", "saved restart population file")
	settingP := flag.String("setting", "", "saved settings file")
	mkenvsP := flag.String("make_envs", "", "Make brand new environments")
	modelP := flag.String("model", "Full", "Model name")
	trajDirP := flag.String("trajdir", "traj", "Directory for trajectory files")
	eStartP := flag.Int("env_start", 0, "starting environment (0, 1, ...)")
	eEndP := flag.Int("env_end", 20, "ending environment")
	seedP := flag.Uint64("seed", 13, "random seed for environments")
	ngenP := flag.Int("ngen", 200, "number of generations per epoch")
	prodP := flag.Bool("production", false, "true if production run")

	flag.Parse()

	var s *multicell.Setting
	if *settingP != "" {
		s = multicell.LoadSetting(*settingP)
	} else {
		s = multicell.GetDefaultSetting()
		s.SetModel(*modelP)
	}
	s.Seed = *seedP
	s.MaxPopulation = *maxpopP
	s.MaxGeneration = *ngenP
	s.Outdir = *trajDirP
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

	if *eStartP >= *eEndP {
		log.Printf("The starting environment must be < the ending environment (%d)", *eEndP)
		os.Exit(1)
	}
	if *eStartP >= len(envs) {
		log.Printf("The starting envirnment must be < %d", len(envs))
		os.Exit(1)
	}

	if *eEndP > len(envs) {
		log.Printf("The ending envirnment must be <= %d", len(envs))
		os.Exit(1)
	}

	var pop multicell.Population
	if *resfileP != "" {
		pop = s.LoadPopulation(*resfileP)
	} else {
		pop = s.NewPopulation(envs[*eStartP])
	}
	return Simulation{
		Setting: s,
		Pop:     pop,
		Envs:    envs,
		Estart:  *eStartP,
		Eend:    *eEndP}

}

func main() {
	t0 := time.Now()
	sim := GetSetting()
	pop := sim.Pop
	sim.Setting.Dump()
	log.Println("pop size: ", len(pop.Indivs))
	var dumpfile string
	for iepoch := sim.Estart; iepoch < sim.Eend; iepoch++ {
		pop.Iepoch = iepoch
		pop, dumpfile = pop.Evolve(sim.Setting, sim.Envs[iepoch])
	}
	log.Printf("Total Time: %v; Dumpifle: %s\n", time.Since(t0), dumpfile)
}

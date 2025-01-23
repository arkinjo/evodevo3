package main

// Genotype-Phenotype Plot

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
	Iepoch  int
	Files   []string // trajectory files
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	envsfileP := flag.String("envs", "", "saved environments JSON file")
	ienvP := flag.Int("ienv", 1, "index of the environment")

	flag.Parse()

	if *settingP == "" {
		log.Fatal("specify a settings file with -setting")
	}
	s := multicell.LoadSetting(*settingP)
	s.Outdir = "gpplot"
	s.Basename += fmt.Sprintf("_ep%2.2d", *ienvP)
	s.ProductionRun = true

	if *envsfileP == "" {
		log.Fatal("specify an environments file with -envs")
	}
	envs := s.LoadEnvs(*envsfileP)

	return Simulation{
		Setting: s,
		Envs:    envs,
		Iepoch:  *ienvP,
		Files:   flag.Args()}

}

func main() {
	t0 := time.Now()
	sim := GetSetting()

	iepoch := sim.Iepoch
	env := sim.Envs[iepoch]

	pop0 := sim.Setting.LoadPopulation(sim.Files[0])
	pop1 := sim.Setting.LoadPopulation(sim.Files[len(sim.Files)-1])

	env0 := sim.Envs[pop0.Iepoch-1]
	env1 := sim.Envs[pop0.Iepoch]

	g0, gaxis := sim.Setting.GetGenomeAxis(pop0, pop1)
	p0, paxis := sim.Setting.GetSelectedPhenoAxis(pop0, pop1, env0, env1)

	log.Printf("Plotting %s epoch %d population under env %d\n",
		sim.Setting.Basename, pop0.Iepoch, iepoch)
	for _, traj := range sim.Files {
		pop := sim.Setting.LoadPopulation(traj)
		if iepoch != pop.Iepoch {
			pop.Initialize(sim.Setting, env)
			pop.Develop(sim.Setting, env)
		}
		pop.GenoPhenoPlot(sim.Setting, p0, paxis, g0, gaxis, env0)
	}

	log.Println("Time: ", time.Since(t0))
}

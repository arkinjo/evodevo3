package main

// Genotype-Phenotype Plot

import (
	"flag"
	"fmt"
	"log"
	//	"sync"
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
	selenv := env.SelectingEnv(sim.Setting)

	pop0 := sim.Setting.LoadPopulation(sim.Files[0])
	pop1 := sim.Setting.LoadPopulation(sim.Files[len(sim.Files)-1])

	env0 := sim.Envs[pop0.Iepoch-1]
	env1 := sim.Envs[pop0.Iepoch]

	g0, gaxis := sim.Setting.GetGenomeAxis(pop0, pop1)
	c0, caxis := sim.Setting.GetCueAxis(env0, env1)
	p0, paxis := sim.Setting.GetPhenoAxis(env0, env1)

	log.Printf("Plotting %s epoch %d population under env %d\n",
		sim.Setting.Basename, pop0.Iepoch, iepoch)

	for _, traj := range sim.Files {
		pop := sim.Setting.LoadPopulation(traj)
		if iepoch != pop.Iepoch {
			pop.Initialize(sim.Setting, env)
			pop.Develop(sim.Setting, selenv)
			pop.Sort()
		}
		ofile := sim.Setting.TrajectoryFilename(pop.Iepoch, pop.Igen, "gpplot")
		pop.Project(sim.Setting, ofile, p0, paxis, g0, gaxis, c0, caxis)
	}

	log.Println("Time: ", time.Since(t0))
}

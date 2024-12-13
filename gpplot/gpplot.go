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

	pop0 := sim.Setting.LoadPopulation(sim.Files[0])
	pop1 := sim.Setting.LoadPopulation(sim.Files[len(sim.Files)-1])

	env0 := sim.Envs[pop0.Iepoch-1]
	env1 := sim.Envs[pop0.Iepoch]
	g0 := multicell.AverageVecs(pop0.GenomeVecs(sim.Setting))
	p0 := env0.SelectingEnv(sim.Setting)
	gaxis := sim.Setting.GetGenomeAxis(&pop0, &pop1)
	paxis := sim.Setting.GetPhenoAxis(env0, env1)

	iepoch := sim.Iepoch
	env := sim.Envs[iepoch]
	selenv := env.SelectingEnv(sim.Setting)
	log.Printf("Plotting %s epoch %d population under env %d\n",
		sim.Setting.Basename, pop0.Iepoch, iepoch)
	for _, traj := range sim.Files {
		pop := sim.Setting.LoadPopulation(traj)
		if iepoch != sim.Iepoch {
			pop.Initialize(sim.Setting, env)
			pop.Develop(sim.Setting, selenv)
		}
		ofile := sim.Setting.TrajectoryFilename(pop.Iepoch, pop.Igen, "gpplot")
		pop.ProjectGenoPheno(sim.Setting, ofile, g0, gaxis, p0, paxis)
	}
	log.Println("Time: ", time.Since(t0))
}

package main

// Develop a population under an arbitrary environment,
// where the population has evolved under a different environment
// whose trajectory files are available.
// Environments must be also saved beforehand.

import (
	"flag"
	"log"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
	Iepoch  int
	Files   []string
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	envsfileP := flag.String("envs", "", "saved environments JSON file")
	ienvP := flag.Int("ienv", 1, "index of the (ancestral) environment")
	flag.Parse()

	if *settingP == "" {
		log.Fatal("specify a settings file with -setting")
	}
	s := multicell.LoadSetting(*settingP)
	s.Basename += "_Anc"
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
	env := sim.Envs[sim.Iepoch]
	selenv := env.SelectingEnv(sim.Setting)
	for _, traj := range sim.Files {
		pop := sim.Setting.LoadPopulation(traj, env)
		pop.Develop(sim.Setting, selenv)
		stats := pop.GetPopStats()
		stats.Print(pop.Iepoch, pop.Igen)
		pop.Dump(sim.Setting)
	}
	log.Println("Time: ", time.Since(t0))
}

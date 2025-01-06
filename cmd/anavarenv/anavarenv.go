package main

// Comparison of plastic responses to various environmental changes.

import (
	"flag"
	"log"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
	Nenvs   int
	Files   []string // trajectory files
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	envsfileP := flag.String("envs", "", "saved Environments JSON file")
	nsamp := flag.Int("n", 10, "number of novel environments")
	flag.Parse()

	if *settingP == "" {
		log.Fatal("specify a settings file with -setting")
	}
	s := multicell.LoadSetting(*settingP)
	var envs []multicell.Environment
	if *envsfileP != "" {
		envs = s.LoadEnvs(*envsfileP)
	} else {
		log.Printf("specify environment file with -envs")
		panic("envs")
	}

	return Simulation{
		Setting: s,
		Envs:    envs,
		Nenvs:   *nsamp,
		Files:   flag.Args()}

}

func main() {
	t0 := time.Now()
	sim := GetSetting()

	for _, traj := range sim.Files {
		pop := sim.Setting.LoadPopulation(traj)
		env0 := sim.Envs[pop.Iepoch]
		pop.AnalyzeVarEnvs(sim.Setting, env0, sim.Nenvs, true)
	}

	log.Println("Time: ", time.Since(t0))
}

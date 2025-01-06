package main

// Comparing the genomes responsible for the adaptive plastic response
// and evolutionary adaptation

import (
	"flag"
	"log"
	"time"

	"github.com/arkinjo/evodevo3/multicell"
)

type Simulation struct {
	Setting *multicell.Setting
	Envs    []multicell.Environment
	Files   []string // trajectory files
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	envsfileP := flag.String("envs", "", "saved Environments JSON file")
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
		Files:   flag.Args()}

}

func main() {
	t0 := time.Now()
	sim := GetSetting()
	// population adapted to the ancestral environment
	pop0 := sim.Setting.LoadPopulation(sim.Files[0])
	// population adapted to the novel environment
	pop1 := sim.Setting.LoadPopulation(sim.Files[len(sim.Files)-1])

	// ancestral environment
	env0 := sim.Envs[pop0.Iepoch-1]
	// novel environment
	env1 := sim.Envs[pop0.Iepoch]

	sim.Setting.AnalyzeAPRGeno(env0, env1, pop0, pop1)
	log.Println("Time: ", time.Since(t0))
}

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
	Igen    int
	Files   []string // trajectory files
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	envsfileP := flag.String("envs", "", "saved environments JSON file")
	ienvP := flag.Int("ienv", 1, "index of the environment")
	igenP := flag.Int("igen", 0, "generation to analyze")
	flag.Parse()

	if *settingP == "" {
		log.Fatal("specify a settings file with -setting")
	}
	s := multicell.LoadSetting(*settingP)
	s.Basename += fmt.Sprintf("_ep%2.2d", *ienvP)
	s.Outdir = "xpca"

	if *envsfileP == "" {
		log.Fatal("specify an environments file with -envs")
	}
	envs := s.LoadEnvs(*envsfileP)

	return Simulation{
		Setting: s,
		Envs:    envs,
		Iepoch:  *ienvP,
		Igen:    *igenP,
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
	p0, paxis := sim.Setting.GetPhenoAxis(pop0, pop1, env0, env1)
	c0, caxis := sim.Setting.GetCueAxis(env0, env1)

	log.Printf("Plotting %s epoch %d population under env %d\n",
		sim.Setting.Basename, pop0.Iepoch, iepoch)
	traj := sim.Files[sim.Igen]
	pop := sim.Setting.LoadPopulation(traj)
	if iepoch != pop.Iepoch {
		pop.Initialize(sim.Setting, env)
		pop.Develop(sim.Setting, env)
	}

	pop.SVDProject(sim.Setting, p0, paxis, g0, gaxis, c0, caxis)

	log.Println("Time: ", time.Since(t0))
}

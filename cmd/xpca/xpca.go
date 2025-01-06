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
	Setting  *multicell.Setting
	Envs     []multicell.Environment
	Iepoch   int
	Selected bool
	Files    []string // trajectory files
}

func GetSetting() Simulation {
	settingP := flag.String("setting", "", "saved settings file")
	selectedP := flag.Bool("selected", false, "saved settings file")
	envsfileP := flag.String("envs", "", "saved environments JSON file")
	ienvP := flag.Int("ienv", 1, "index of the environment")
	flag.Parse()

	if *settingP == "" {
		log.Fatal("specify a settings file with -setting")
	}
	s := multicell.LoadSetting(*settingP)
	s.Basename += fmt.Sprintf("_ep%2.2d", *ienvP)

	if *envsfileP == "" {
		log.Fatal("specify an environments file with -envs")
	}
	envs := s.LoadEnvs(*envsfileP)

	return Simulation{
		Setting:  s,
		Envs:     envs,
		Iepoch:   *ienvP,
		Selected: *selectedP,
		Files:    flag.Args()}

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
	p0, paxis := sim.Setting.GetPhenoAxis(sim.Selected, env0, env1)
	log.Printf("Selected phenotype only? %b\n", sim.Selected)
	c0, caxis := sim.Setting.GetCueAxis(env0, env1)

	log.Printf("Plotting %s epoch %d population under env %d\n",
		sim.Setting.Basename, pop0.Iepoch, iepoch)
	traj := sim.Files[0]
	pop := sim.Setting.LoadPopulation(traj)
	if iepoch != pop.Iepoch {
		pop.Initialize(sim.Setting, env)
		pop.Develop(sim.Setting, selenv)
	}

	pop.SVDProject(sim.Setting, sim.Selected, p0, paxis, g0, gaxis, c0, caxis)

	log.Println("Time: ", time.Since(t0))
}

package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/arkinjo/evodevo3/multicell"
	//	"log"
	"os"
	//	"strings"
	//	"math"
)

func dumpjson(js []byte) {
	fmt.Println(len(js))
	os.Stdout.Write(js)
}

func main() {
	s := multicell.GetDefaultSetting()
	s.SetOmega()
	s.MaxPopulation = 100
	s.ProductionRun = false
	s.Outdir = "traj"
	s.FullModel()
	s.SaveEnvs("envs.json", 50)
	envs := s.LoadEnvs("envs.json")
	env := envs[0]
	pop := s.NewPopulation(env)
	fmt.Println("#With_cue= ", s.WithCue)
	pop.Evolve(s, 10, env)
	ofilename := pop.Dump(s)
	pop = s.LoadPopulation(ofilename, env)
	fmt.Println(s.SelectingEnv(env))
	fmt.Println(pop.Indivs[0].SelectedPhenotype())
}

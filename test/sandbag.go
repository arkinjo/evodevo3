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
	s.EnvNoise = 0.0
	s.MaxPopulation = 100
	s.MaxGeneration = 10
	s.ProductionRun = false
	s.SetModel("Full")
	s.Outdir = "traj"
	fmt.Println("#With_cue= ", s.WithCue)
	//	s.SaveEnvs("envs.json", 50)
	envs := s.LoadEnvs("envs.json")
	env := envs[0]
	//	pop := s.NewPopulation(env)
	//pop.Evolve(s, env)
	//ofilename := pop.Dump(s)
	ofilename := "traj/Full_00_010.traj.gz"
	env = envs[1]
	pop := s.LoadPopulation(ofilename, env)
	pop.DumpJSON(s)
	fmt.Println("nindiv= ", len(pop.Indivs))
	fmt.Println(s.SelectingEnv(env))
	fmt.Println(pop.Indivs[0].SelectedPhenotype())
	fmt.Println(len(pop.Indivs[0].Genome.ToVec(s)))
}

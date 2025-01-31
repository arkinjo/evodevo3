package multicell

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

var rng = rand.New(rand.NewPCG(13, 97))

type Environment = Vec

type EnvironmentS []Environment

type CellEnvs struct {
	Tops    []Vec
	Bottoms []Vec
	Rights  []Vec
	Lefts   []Vec
}

func (env Environment) Len() int {
	return len(env)
}

func (s *Setting) NewEnvironment() Environment {
	lenenv := s.LenFace * 4
	env := make([]float64, lenenv)
	for i := range lenenv {
		if i%2 == 0 {
			env[i] = 1
		} else {
			env[i] = -1
		}
	}
	return env
}

func (env Environment) Left(s *Setting) Vec {
	return env[0:s.LenFace]
}

func (env Environment) Top(s *Setting) Vec {
	return env[s.LenFace : s.LenFace*2]
}

func (env Environment) Right(s *Setting) Vec {
	return env[s.LenFace*2 : s.LenFace*3]
}

func (env Environment) Bottom(s *Setting) Vec {
	return env[s.LenFace*3 : s.LenFace*4]
}

func (env Environment) Face(s *Setting, iface int) Vec {
	var v Vec
	switch iface {
	case Left:
		v = env.Left(s)
	case Top:
		v = env.Top(s)
	case Right:
		v = env.Right(s)
	case Bottom:
		v = env.Bottom(s)
	default:
		log.Fatal("(*env).Face: unknown face")
	}

	return v
}

func (env Environment) Compare(env0 Environment) float64 {
	denv := make(Environment, len(env))
	d := denv.Diff(env, env0).Norm1() / 2
	return d
}

func (env Environment) SelectingEnv(s *Setting) Vec {
	return env.Left(s)
}

func (env Environment) AddNoise(p float64) Vec {
	nvec := env.Clone()
	for i := range nvec {
		if rand.Float64() < p {
			nvec[i] *= -1
		}
	}

	return nvec
}

func (env Environment) ChangeEnv(s *Setting) Environment {
	nblk := s.LenFace / s.LenBlock
	nflip := int(s.Denv * float64(nblk))
	nenv := env.Clone()

	for iface := range NumFaces {
		i := iface * s.LenFace
		for _, p := range rng.Perm(nblk)[:nflip] {
			j := i + p*nblk
			for k := range s.LenBlock {
				nenv[j+k] *= -1
			}
		}
	}
	return nenv
}

func (env Environment) BlockFlip(s *Setting, ref Environment) Environment {
	var nenv Environment
	if rng.Float64() < math.Exp(-0.1) {
		return env
	} else {
		nenv = ref.Clone()
	}

	nblk := len(env) / s.LenBlock
	for ib := range nblk {
		i := ib * s.LenBlock
		if rng.Float64() < s.Penv01 {
			for j := range s.LenBlock {
				nenv[i+j] *= -1
			}
		}
	}

	return nenv
}

func (env Environment) MarkovFlip(s *Setting, ref Environment) Environment {
	nenv := env.Clone()
	nblk := len(env) / s.LenBlock
	for ib := range nblk {
		i := ib * s.LenBlock
		r2v := (ref[i] == env[i] && rng.Float64() < s.Penv01)
		v2r := (ref[i] != env[i] && rng.Float64() < s.Penv10)
		if r2v || v2r {
			for j := range s.LenBlock {
				nenv[i+j] *= -1
			}
		}
	}

	return nenv
}

func (env Environment) GenerateEnvs(s *Setting, nepochs int) EnvironmentS {

	envs := make([]Environment, nepochs)
	envs[0] = env
	for n := range nepochs {
		if n == 0 {
			continue
		}
		envs[n] = envs[n-1].ChangeEnv(s)

	}
	return envs
}

func (s *Setting) SaveEnvs(filename string, nepochs int) EnvironmentS {
	env0 := s.NewEnvironment()
	envs := env0.GenerateEnvs(s, nepochs)
	envs.DumpEnvs(filename)
	return envs
}

func (envs EnvironmentS) DumpEnvs(filename string) {
	fout, err := os.Create(filename)
	JustFail(err)
	defer fout.Close()

	fmt.Fprintf(fout, "[")
	for i, env := range envs {
		json, err := json.Marshal(env)
		JustFail(err)
		fmt.Fprintf(fout, "%s", json)
		if i < len(envs)-1 {
			fmt.Fprintf(fout, ",\n")
		}
	}
	fmt.Fprintf(fout, "]\n")
}

func (s *Setting) LoadEnvs(filename string) []Environment {
	var envs []Environment
	buffer, err := os.ReadFile(filename)
	JustFail(err)
	err = json.Unmarshal(buffer, &envs)
	JustFail(err)
	return envs
}

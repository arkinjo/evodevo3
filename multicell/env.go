package multicell

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"os"
)

type Environment = Vec

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
		if i < lenenv/2 {
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

func (env Environment) AddNoise(s *Setting) Environment {
	nenv := make([]float64, env.Len())
	for i, v := range env {
		if rand.Float64() < s.EnvNoise {
			nenv[i] = -v
		} else {
			nenv[i] = v
		}
	}

	return nenv
}

func (env Environment) SelectingEnv(s *Setting) Vec {
	return env.Left(s)
}

func (env Environment) ChangeEnv(s *Setting, rng *rand.Rand) Environment {
	lenv := env.Len()
	nflip := int(s.Denv * float64(lenv))
	nenv := make(Vec, lenv)

	for i, j := range rng.Perm(lenv) {
		nenv[i] = env[i]
		if j < nflip {
			nenv[i] *= -1
		}
	}
	return nenv
}

func (s *Setting) SaveEnvs(filename string, nepochs int) []Environment {
	rng := rand.New(rand.NewPCG(s.Seed, s.Seed+1397))
	env := s.NewEnvironment()
	envs := make([]Environment, nepochs)
	for i := range nepochs {
		env = env.ChangeEnv(s, rng)
		envs[i] = env
	}
	json, err := json.Marshal(envs)
	JustFail(err)
	os.WriteFile(filename, json, 0644)
	return envs
}

func (s *Setting) LoadEnvs(filename string) []Environment {
	var envs []Environment
	buffer, err := os.ReadFile(filename)
	JustFail(err)
	err = json.Unmarshal(buffer, &envs)
	JustFail(err)
	return envs
}

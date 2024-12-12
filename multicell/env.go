package multicell

import (
	"encoding/json"
	"math/rand/v2"
	"os"
)

type Environment struct {
	V []float64
}

type CellEnvs struct {
	Tops    []Vec
	Bottoms []Vec
	Rights  []Vec
	Lefts   []Vec
}

func (env *Environment) LenEnv() int {
	return len(env.V)
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
	return Environment{V: env}
}

func (env *Environment) Left(s *Setting) Vec {
	return env.V[0:s.LenFace]
}

func (env *Environment) Top(s *Setting) Vec {
	return env.V[s.LenFace : s.LenFace*2]
}

func (env *Environment) Right(s *Setting) Vec {
	return env.V[s.LenFace*2 : s.LenFace*3]
}

func (env *Environment) Bottom(s *Setting) Vec {
	return env.V[s.LenFace*3:]
}

func (env *Environment) AddNoise(s *Setting) Environment {
	nenv := make([]float64, env.LenEnv())
	for i, v := range env.V {
		if rand.Float64() < s.EnvNoise {
			nenv[i] = -v
		} else {
			nenv[i] = v
		}
	}

	return Environment{V: nenv}
}

func (env *Environment) SelectingEnv(s *Setting) Vec {
	return env.Left(s)
}

func (env *Environment) ChangeEnv(s *Setting, rng *rand.Rand) Environment {
	lenv := env.LenEnv()
	nflip := int(s.Denv * float64(lenv))
	nenv := make(Vec, lenv)

	for i, j := range rng.Perm(lenv) {
		nenv[i] = env.V[i]
		if j < nflip {
			nenv[i] *= -1
		}
	}
	return Environment{V: nenv}
}

func (s *Setting) SaveEnvs(filename string, nepochs int) {
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
}

func (s *Setting) LoadEnvs(filename string) []Environment {
	var envs []Environment
	buffer, err := os.ReadFile(filename)
	JustFail(err)
	err = json.Unmarshal(buffer, &envs)
	JustFail(err)
	return envs
}

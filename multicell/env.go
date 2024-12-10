package multicell

import (
	"encoding/json"
	"math/rand/v2"
	"os"
)

type Environment []float64

type CellEnvs struct {
	Tops    []Vec
	Bottoms []Vec
	Rights  []Vec
	Lefts   []Vec
}

func (s *Setting) LenEnv() int {
	return s.LenFace * (s.NumCellX + s.NumCellY) * 2 * NumFaces
}

func (s *Setting) NewEnvironment() Environment {
	lenenv := s.LenEnv()
	env := make([]float64, lenenv)
	for i, j := range rand.Perm(lenenv) {
		if i < lenenv/2 {
			env[j] = 1
		} else {
			env[j] = -1
		}
	}
	return env
}

func (s *Setting) AddNoise(env Environment) Environment {
	nenv := make([]float64, s.LenEnv())
	for i, v := range env {
		if rand.Float64() < s.EnvNoise {
			nenv[i] = -v
		} else {
			nenv[i] = v
		}
	}

	return nenv
}

func (s *Setting) NewCellEnvs(env Environment) CellEnvs {
	var env0 Vec
	if s.WithCue {
		env0 = env
	} else {
		env0 = NewVec(s.LenEnv(), 1.0)
	}
	nenv := s.AddNoise(env0)
	lenv := s.LenFace
	lenx := lenv * s.NumCellX
	leny := lenv * s.NumCellY

	left := nenv[0:leny]
	top := nenv[leny : leny+lenx]
	right := nenv[leny+lenx : leny*2+lenx]
	bottom := nenv[leny*2:]

	lenx1 := lenx / s.NumCellX
	leny1 := leny / s.NumCellY

	lefts := make([]Vec, s.NumCellY)
	rights := make([]Vec, s.NumCellY)
	tops := make([]Vec, s.NumCellX)
	bottoms := make([]Vec, s.NumCellX)

	for i := 0; i < s.NumCellX; i++ {
		tops[i] = top[i*lenx1 : (i+1)*lenx1]
		bottoms[i] = bottom[i*lenx1 : (i+1)*lenx1]
	}

	for i := 0; i < s.NumCellY; i++ {
		lefts[i] = left[i*leny1 : (i+1)*leny1]
		rights[i] = right[i*leny1 : (i+1)*leny1]
	}

	return CellEnvs{
		Tops:    tops,
		Bottoms: bottoms,
		Rights:  rights,
		Lefts:   lefts,
	}
}

func (s *Setting) SelectingEnv(env Environment) Environment {
	return env[0 : s.LenFace*s.NumCellY]
}

func (s *Setting) ChangeEnv(env Environment) Environment {
	ndenv := int(s.Denv * float64(s.LenEnv()))
	nenv := make(Environment, s.LenEnv())
	copy(nenv, env)

	indices := make([]int, s.LenEnv())
	for i := range indices {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	for _, i := range indices[:ndenv] {
		nenv[i] *= -1
	}
	return nenv
}

func (s *Setting) SaveEnvs(filename string, nepochs int) {
	env := s.NewEnvironment()
	envs := make([]Environment, nepochs)
	for i := range nepochs {
		env = s.ChangeEnv(env)
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

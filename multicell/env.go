package multicell

import (
	"math/rand/v2"
)

type Environment []float64

type Cell_envs struct {
	Tops    []Vec
	Bottoms []Vec
	Rights  []Vec
	Lefts   []Vec
}

func (s *Setting) NewEnvironment() Environment {
	env := make([]float64, s.Num_env)
	for i, j := range rand.Perm(s.Num_env) {
		if i < s.Num_env/2 {
			env[j] = 1
		} else {
			env[j] = -1
		}
	}
	return env
}

func (s *Setting) Add_noise(env Environment) Environment {
	nenv := make([]float64, s.Num_env)
	for i, v := range env {
		if rand.Float64() < s.Env_noise {
			nenv[i] = -v
		} else {
			nenv[i] = v
		}
	}

	return nenv
}

func (s *Setting) NewCell_envs(env Environment) Cell_envs {
	var env0 Vec
	if s.With_cue {
		env0 = env
	} else {
		env0 = NewVec(s.Num_env, 1.0)
	}
	nenv := s.Add_noise(env0)
	lenv := s.Num_components[0] / 4
	lenx := lenv * s.Num_cell_x
	leny := lenv * s.Num_cell_y

	left := nenv[0:leny]
	top := nenv[leny : leny+lenx]
	right := nenv[leny+lenx : leny*2+lenx]
	bottom := nenv[leny*2:]

	lenx1 := lenx / s.Num_cell_x
	leny1 := leny / s.Num_cell_x

	lefts := make([]Vec, s.Num_cell_y)
	rights := make([]Vec, s.Num_cell_y)
	tops := make([]Vec, s.Num_cell_x)
	bottoms := make([]Vec, s.Num_cell_x)

	for i := 0; i < s.Num_cell_x; i++ {
		tops[i] = top[i*lenx1 : (i+1)*lenx1]
		bottoms[i] = bottom[i*lenx1 : (i+1)*lenx1]
	}

	for i := 0; i < s.Num_cell_y; i++ {
		lefts[i] = left[i*leny1 : (i+1)*leny1]
		rights[i] = right[i*leny1 : (i+1)*leny1]
	}

	return Cell_envs{
		Tops:    tops,
		Bottoms: bottoms,
		Rights:  rights,
		Lefts:   lefts,
	}
}

func (s *Setting) Selecting_env(env Environment) Environment {
	return env[0 : s.Num_components[0]*s.Num_cell_y/4]
}

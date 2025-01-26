package multicell_test

import (
	"fmt"
	"testing"

	"github.com/arkinjo/evodevo3/multicell"
)

func TestEnvChange(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Denv = 0.5
	envs := s.SaveEnvs(ENVSFILE, 50)
	nexp := int(s.Denv * float64(len(envs[0])))

	for n, env := range envs {
		if n == 0 {
			continue
		}
		env0 := envs[n-1]
		diff := make(multicell.Vec, len(env))
		ndiff := int(diff.Diff(env0, env).ScaleBy(0.5).Norm1())
		if ndiff != nexp {
			t.Errorf("EnvChange: %d; expected %d\n", ndiff, nexp)
		}
	}
}

func TestEnvCue(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	envs := s.SaveEnvs(ENVSFILE, 50)
	s.NumBlocks = 2
	s.LenBlock = 5
	env := envs[1]
	cue := env.GetCue(s)
	ndiff := 0
	for i, v := range cue {
		if v != env[i] {
			ndiff += 1
		}
	}
	for i, e := range env {
		if e != cue[i] {
			fmt.Printf("env/cue: %d %2.0f %2.0f\n", i, e, cue[i])
		}
	}
	nexp := s.NumBlocks * s.LenBlock
	if ndiff != nexp {
		t.Errorf("Cue difference %d; expected %d\n", ndiff, nexp)
	}
}

func TestEnvNoCue(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	s.Outdir = "traj"
	s.WithCue = false
	envs := s.SaveEnvs(ENVSFILE, 50)
	s.NumBlocks = 2
	s.LenBlock = 5
	env := envs[1]
	cue := env.GetCue(s)
	ndiff := 0
	for _, v := range cue {
		if v != -1.0 {
			ndiff += 1
		}
	}

	nexp := s.NumBlocks * s.LenBlock
	if ndiff != nexp {
		t.Errorf("Cue difference %d; expected %d\n", ndiff, nexp)
	}
}

func TestEnvironment(t *testing.T) {
	s := multicell.GetDefaultSetting("Full")
	envs := s.SaveEnvs(ENVSFILE, 50)
	if len(envs) != 50 {
		t.Errorf("len(envs)= %d; want 50", len(envs))
	}
	for i, env := range envs {
		if env.Len() != s.LenFace*4 {
			t.Errorf("env[%d].Len()= %d; want %d", i, env.Len(), s.LenFace*4)
		}
	}
}

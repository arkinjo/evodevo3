package evodevo3_test

import (
	"github.com/arkinjo/evodevo3/multicell"
	"log"
	"testing"
)

func TestSetting(t *testing.T) {
	log.Println("Setting")
	s := multicell.GetDefaultSetting()
	models := []string{"Full", "NoCue", "NoHie", "NoDev", "Null",
		"NullCue", "NullHie", "NullDev"}
	for _, model := range models {
		s.SetModel(model)
		got := len(s.Topology)
		if got != s.NumLayers {
			t.Errorf("got len(Topology) = %d; want %d", got, s.NumLayers)
		}
	}
}

func TestEnvironment(t *testing.T) {
	s := multicell.GetDefaultSetting()
	s.SaveEnvs("envs.json", 50)
	envs := s.LoadEnvs("envs.json")
	if len(envs) != 50 {
		t.Errorf("len(envs)= %d; want 50", len(envs))
	}
	for i, env := range envs {
		if env.Len() != s.LenFace*4 {
			t.Errorf("env[%d].Len()= %d; want %d", i, env.Len(), s.LenFace*4)
		}
	}
}

package config

import (
	"os"
	"github.com/goccy/go-yaml"
)

type EBPFConfig struct {
	EBPFProgram struct {
		Name             string `yaml:"name"`
		Type             string `yaml:"type"`
		AttachPoint      string `yaml:"attach_point"`
		KernelVersion    string `yaml:"kernel_version"`
		VerifierOptions  map[string]interface{} `yaml:"verifier_options"`
		Maps             []map[string]interface{} `yaml:"maps"`
		Features         []string `yaml:"features"`
		Dependencies     map[string]string `yaml:"dependencies"`
		Loading          map[string]bool `yaml:"loading"`
		Testing          map[string]interface{} `yaml:"testing"`
	} `yaml:"ebpf_program"`
}

func ParseConfig(path string) (*EBPFConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg EBPFConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}


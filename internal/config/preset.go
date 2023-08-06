package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Preset struct {
	Name         string             `yaml:"name" validate:"required"`
	Namespace    string             `yaml:"namespace" validate:"required"`
	Deployments  []ScalableResource `yaml:"deployments" validate:"required,dive"`
	StatefulSets []ScalableResource `yaml:"statefulSets" validate:"required,dive"`
}

func LoadPreset(filepath string) (*Preset, error) {
	var bytes, err = os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var preset Preset
	return &preset, yaml.Unmarshal(bytes, &preset)
}

func (p *Preset) SavePreset(filepath string) error {
	var bytes, err = yaml.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, bytes, 0644)
}

func (p *Preset) Validate() {

}

package config

type ScalableResource struct {
	Name     string `yaml:"name" validate:"required"`
	Replicas int32  `yaml:"replicas" validate:"required,gte=0"`
}

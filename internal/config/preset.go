package config

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/th3shadowbroker/legion/internal/kube"
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

	if _, err := os.Stat(filepath); err == nil {
		pterm.Error.Printfln("Error: File already exists")
		return nil
	}

	return os.WriteFile(filepath, bytes, 0644)
}

func (p *Preset) Apply(client *kube.KubeClient) {
	pterm.Info.Printfln("Applying preset %s to namespace %s", p.Name, p.Namespace)
	p.ApplyDeployments(client)
	p.ApplyStatefulSet(client)
}

func (p *Preset) ApplyDeployments(client *kube.KubeClient) {
	if len(p.Deployments) > 0 {
		var progress, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Applying deployments to namespace %s", p.Namespace))
		for _, deployment := range p.Deployments {
			progress.UpdateText(fmt.Sprintf("Scaling %s to %d", deployment.Name, deployment.Replicas))
			_, exists, err := client.SetDeploymentScale(p.Namespace, deployment.Name, deployment.Replicas)
			if err != nil {
				pterm.Error.Println("Could not scale deployment %s in namespace %s: %s", deployment.Name, p.Namespace, err.Error())
				continue
			}

			if !exists {
				pterm.Warning.Println("Deployment %s not found in namespace %s", deployment.Name, p.Namespace)
				continue
			}

			pterm.Success.Printfln("Scaled deployment %s to %d", deployment.Name, deployment.Replicas)
		}
		progress.Info(fmt.Sprintf("Processed %d deployments", len(p.Deployments)))
	}
}

func (p *Preset) ApplyStatefulSet(client *kube.KubeClient) {
	if len(p.StatefulSets) > 0 {
		var progress, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Applying statefulsets to namespace %s", p.Namespace))
		for _, statefulSet := range p.StatefulSets {
			progress.UpdateText(fmt.Sprintf("Scaling %s to %d", statefulSet.Name, statefulSet.Replicas))
			_, exists, err := client.SetDeploymentScale(p.Namespace, statefulSet.Name, statefulSet.Replicas)
			if err != nil {
				pterm.Error.Println("Could not scale statefulset %s in namespace %s: %s", statefulSet.Name, p.Namespace, err.Error())
				continue
			}

			if !exists {
				pterm.Warning.Println("Statefulset %s not found in namespace %s", statefulSet.Name, p.Namespace)
				continue
			}

			pterm.Success.Printfln("Scaled statefulset %s to %d", statefulSet.Name, statefulSet.Replicas)
		}
		progress.Info(fmt.Sprintf("Processed %d statefulsets", len(p.StatefulSets)))
	}
}

func (p *Preset) Populate(client *kube.KubeClient) {
	pterm.Info.Printfln("Creating preset %s from namespace %s", p.Name, p.Namespace)
	p.PopulateDeployments(client)
	p.PopulateStatefulSets(client)
}

func (p *Preset) PopulateDeployments(client *kube.KubeClient) {
	spinner, _ := pterm.DefaultSpinner.Start("Fetching deployments...")
	deployments, err := client.GetDeploymentsInNamespace(p.Namespace)
	if err != nil {
		pterm.Error.Printfln("Could not fetch deployments: %s", err.Error())
	}

	for _, deployment := range deployments.Items {
		p.Deployments = append(p.Deployments, ScalableResource{
			Name:     deployment.Name,
			Replicas: *deployment.Spec.Replicas,
		})
		pterm.Success.Printfln("Added %s", deployment.Name)
	}
	spinner.Success(fmt.Sprintf("Processed %d deployments", len(deployments.Items)))
}

func (p *Preset) PopulateStatefulSets(client *kube.KubeClient) {
	spinner, _ := pterm.DefaultSpinner.Start("Fetching statefulsets...")
	statefulSets, err := client.GetStatefulSetsInNamespace(p.Namespace)
	if err != nil {
		pterm.Error.Printfln("Could not fetch statefulsets: %s", err.Error())
	}

	for _, deployment := range statefulSets.Items {
		p.Deployments = append(p.Deployments, ScalableResource{
			Name:     deployment.Name,
			Replicas: *deployment.Spec.Replicas,
		})
		pterm.Success.Printfln("Added %s", deployment.Name)
	}
	spinner.Success(fmt.Sprintf("Processed %d statefulsets", len(statefulSets.Items)))
}

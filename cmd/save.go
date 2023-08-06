package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/th3shadowbroker/legion/internal/config"
	"github.com/th3shadowbroker/legion/internal/kube"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Saves the current replica configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		var client, err = kube.NewKubeClientFromCommand(cmd)
		if err != nil {
			return err
		}

		var preset = new(config.Preset)
		presetName, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")

		// Metadata
		preset.Name = presetName
		preset.Namespace = namespace

		// Deployments
		var spinner, _ = pterm.DefaultSpinner.Start("Fetching deployments...")
		deployments, err := client.GetDeploymentsInNamespace(namespace)
		if err != nil {
			spinner.Fail("Could not retrieve deployments:", err.Error())
		}

		for _, deployment := range deployments.Items {
			var deploymentResource = config.ScalableResource{
				Name:     deployment.Name,
				Replicas: *deployment.Spec.Replicas,
			}
			preset.Deployments = append(preset.Deployments, deploymentResource)
		}
		spinner.Success("Fetched ", len(preset.Deployments), " deployments")

		// Statefulsets
		spinner, _ = pterm.DefaultSpinner.Start("Fetching statefulsets...")
		statefulSets, err := client.GetStatefulSetsInNamespace(namespace)
		if err != nil {
			spinner.Fail("Could not retrieve statefulsets:", err.Error())
		}

		for _, statefulset := range statefulSets.Items {
			var statefulsetResource = config.ScalableResource{
				Name:     statefulset.Name,
				Replicas: *statefulset.Spec.Replicas,
			}
			preset.StatefulSets = append(preset.StatefulSets, statefulsetResource)
		}
		spinner.Success("Fetched ", len(preset.StatefulSets), " statefulsets")

		var path = filepath.Join(".", presetName+".yml")
		if err := preset.SavePreset(path); err != nil {
			pterm.Fatal.Println("Could not write preset:", err.Error())
		}

		return nil
	},
}

func init() {
	var defaultKubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	if kubeConfigEnv, ok := os.LookupEnv("KUBECONFIG"); ok {
		defaultKubeConfig = kubeConfigEnv
	}

	saveCmd.PersistentFlags().StringP("kubeconfig", "k", defaultKubeConfig, "pass the configuration that should be used")
	saveCmd.PersistentFlags().BoolP("service-account", "s", false, "auto-configure via a mounted service account")
	saveCmd.PersistentFlags().BoolP("dry-run", "d", false, "let's pretend that you know what you're doing")
	saveCmd.MarkFlagsMutuallyExclusive("kubeconfig", "service-account")

	saveCmd.Flags().StringP("name", "p", "", "set the preset name")
	saveCmd.Flags().StringP("namespace", "n", "", "set the namespace to read from")

	_ = saveCmd.MarkFlagRequired("name")
	_ = saveCmd.MarkFlagRequired("namespace")
}

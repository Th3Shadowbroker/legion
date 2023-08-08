package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "legion",
	Short: "A tool for preset-based scaling of kubernetes deployments",
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(restoreCmd)

	var defaultKubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	if kubeConfigEnv, ok := os.LookupEnv("KUBECONFIG"); ok {
		defaultKubeConfig = kubeConfigEnv
	}

	rootCmd.PersistentFlags().StringP("kubeconfig", "k", defaultKubeConfig, "pass the configuration that should be used")
	rootCmd.PersistentFlags().BoolP("service-account", "s", false, "auto-configure via a mounted service account")
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "let's pretend that you know what you're doing")
	rootCmd.MarkFlagsMutuallyExclusive("kubeconfig", "service-account")
}

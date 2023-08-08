package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/th3shadowbroker/legion/internal/config"
	"github.com/th3shadowbroker/legion/internal/kube"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores a replica configuration",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kube.NewKubeClientFromCommand(cmd.Parent())
		if err != nil {
			pterm.Error.Printfln("Could not connect to cluster: %s", err.Error())
			return
		}

		presetPath, _ := cmd.Flags().GetString("file")
		preset, err := config.LoadPreset(presetPath)
		if err != nil {
			pterm.Error.Printfln("Could not read config: %s", err.Error())
			return
		}
		preset.Apply(client)
	},
}

func init() {
	restoreCmd.Flags().StringP("file", "f", "", "set the preset file to read from")
	_ = restoreCmd.MarkFlagFilename("file")
	_ = restoreCmd.MarkFlagRequired("file")
}

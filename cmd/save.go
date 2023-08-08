package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/th3shadowbroker/legion/internal/config"
	"github.com/th3shadowbroker/legion/internal/kube"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Saves the current replica configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := kube.NewKubeClientFromCommand(cmd.Parent())
		if err != nil {
			pterm.Error.Printfln("Could not connect to cluster: %s", err.Error())
			return nil
		}

		var preset = new(config.Preset)
		presetName, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")

		preset.Name = presetName
		preset.Namespace = namespace
		preset.Populate(client)

		presetFilePath, _ := cmd.Flags().GetString("file")
		if presetFilePath == "" {
			presetFilePath = fmt.Sprintf("%s.yml", presetName)
		}

		return preset.SavePreset(presetName)
	},
}

func init() {
	saveCmd.Flags().StringP("file", "f", "", "the name of the file to write to (<preset-name>.yml if not set)")
	saveCmd.Flags().StringP("name", "p", "", "set the preset name")
	saveCmd.Flags().StringP("namespace", "n", "", "set the namespace to read from")

	_ = saveCmd.MarkFlagRequired("name")
	_ = saveCmd.MarkFlagRequired("namespace")
}

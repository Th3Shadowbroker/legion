package cmd

import "github.com/pterm/pterm"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Printfln("Error: %s", err.Error())
	}
}

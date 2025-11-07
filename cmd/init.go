/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"projman/internal/config"
	"projman/internal/project"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Seleziona i progetti da includere nella gestione",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Controlla che sia stata passata la directory root
		if len(args) < 1 {
			pterm.Error.Println("Devi specificare la directory root dei progetti")
			return nil
		}

		// Controlla che la directory esista
		root, err := config.CheckAndGetDirectory(args[0])
		if err != nil {
			pterm.Error.Println("Immetti una directory valida")
			return nil
		}

		// Scansione dei progetti Maven
		pterm.Info.Println("Scansione dei progetti in " + root)
		projs, _ := project.Discover(root)
		if len(projs) == 0 {
			pterm.Warning.Println("Nessun progetto Maven trovato")
			return nil
		}

		// Prompt per la selezione dei progetti
		names := project.Names(projs)
		var selected []string
		prompt := &survey.MultiSelect{
			Message:  "Seleziona i progetti da includere:",
			Options:  names,
			Default:  names,
			PageSize: 10,
		}

		if err := survey.AskOne(prompt, &selected); err != nil {
			return err
		}

		// Salva lo config in json
		return config.SaveSettings(config.Config{SelectedProjects: selected})

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

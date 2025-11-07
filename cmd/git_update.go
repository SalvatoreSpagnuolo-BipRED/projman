/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/executil"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Git Update")

		// Carica la configurazione
		cfg, err := config.LoadSettings()
		if err != nil {
			pterm.Error.Println("Errore nel caricamento della configurazione:", err)
			return
		}

		// Esegui git pull per ogni progetto selezionato
		spinner, _ := pterm.DefaultSpinner.Start("Eseguo git pull")
		for _, p := range cfg.SelectedProjects {
			spinner.UpdateText("Aggiorno il progetto: " + p)
			projectPath := filepath.Join(cfg.RootOfProjects, p)
			if err := executil.Run("git", "-C", projectPath, "pull"); err != nil {
				pterm.Error.Println("Errore durante l'aggiornamento del progetto " + p + ": " + err.Error())
			} else {
				pterm.Success.Println("Progetto aggiornato: " + p)
			}
		}

		spinner.Success("Git pull completato")

	},
}

func init() {
	gitCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

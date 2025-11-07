/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/executil"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var runTests bool

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Esegue mvn install su tutti i progetti selezionati",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Maven Install")

		// Carica la configurazione
		cfg, err := config.LoadSettings()
		if err != nil {
			pterm.Error.Println("Errore nel caricamento della configurazione:", err)
			return
		}

		pterm.Info.Println("test abilitati")

		// Esegui git pull per ogni progetto selezionato
		spinner, _ := pterm.DefaultSpinner.Start("Eseguo mvn install")
		for _, p := range cfg.SelectedProjects {
			spinner.UpdateText("Installo il progetto: " + p)
			pomPath := filepath.Join(cfg.RootOfProjects, p, "pom.xml")
			args := []string{"-f", pomPath, "install"}
			if !runTests {
				args = append(args, "-DskipTests=true")
			}

			if err := executil.Run("mvn", args...); err != nil {
				pterm.Error.Println("Errore durante install del progetto " + p + ": " + err.Error())
			} else {
				pterm.Success.Println("Progetto installato: " + p)
			}
		}

		spinner.Success("mvn install completato")

	},
}

func init() {
	mvnCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&runTests, "tests", "t", false, "Abilita i test")
}

package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Gestisce le operazioni Git sui progetti selezionati",
	Long: `Esegue operazioni Git su tutti i progetti precedentemente selezionati.
Supporta vari comandi Git come update (pull) per mantenere i progetti aggiornati.`,

	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Git Command")
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}

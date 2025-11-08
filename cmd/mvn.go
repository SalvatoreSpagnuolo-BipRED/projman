package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var mvnCmd = &cobra.Command{
	Use:   "mvn",
	Short: "Esegue comandi Maven sui progetti selezionati",
	Long: `Permette di eseguire comandi Maven su tutti i progetti selezionati.
Questo comando supporta varie operazioni Maven come install, clean, compile, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Maven Command")
	},
}

func init() {
	rootCmd.AddCommand(mvnCmd)
}

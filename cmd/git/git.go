package git

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// gitCmd rappresenta il comando parent per tutte le operazioni Git
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Gestisce le operazioni Git sui progetti selezionati",
	Long: `Esegue operazioni Git su tutti i progetti precedentemente selezionati.
Supporta vari comandi Git come update (pull) per mantenere i progetti aggiornati.

Per utilizzare questo comando, è necessario specificare un sottocomando.
Esempi:
  projman git update    - Aggiorna tutti i progetti con git pull/merge`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Error.Println("Devi specificare un sottocomando (es. 'update')")
		pterm.Info.Println("Usa 'projman git --help' per vedere i sottocomandi disponibili")
		_ = cmd.Help()
	},
}

// GetGitCmd restituisce il comando git per la registrazione in root
func GetGitCmd() *cobra.Command {
	return gitCmd
}

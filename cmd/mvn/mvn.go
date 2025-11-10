package mvn

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// mvnCmd rappresenta il comando parent per tutte le operazioni Maven
var mvnCmd = &cobra.Command{
	Use:   "mvn",
	Short: "Esegue comandi Maven sui progetti selezionati",
	Long: `Permette di eseguire comandi Maven su tutti i progetti selezionati.
Questo comando supporta varie operazioni Maven come install, clean, compile, etc.

Per utilizzare questo comando, è necessario specificare un sottocomando.
Esempi:
  projman mvn install         - Esegue mvn install senza test
  projman mvn install --tests - Esegue mvn install con test abilitati`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Error.Println("Devi specificare un sottocomando (es. 'install')")
		pterm.Info.Println("Usa 'projman mvn --help' per vedere i sottocomandi disponibili")
		_ = cmd.Help()
	},
}

// GetMvnCmd restituisce il comando mvn per la registrazione in root
func GetMvnCmd() *cobra.Command {
	return mvnCmd
}

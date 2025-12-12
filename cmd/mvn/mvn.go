package mvn

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/cmdutil"
	"github.com/spf13/cobra"
)

// MvnCmd rappresenta il comando parent per tutte le operazioni Maven
var MvnCmd = &cobra.Command{
	Use:   "mvn",
	Short: "Esegue comandi Maven sui progetti selezionati",
	Long: `Permette di eseguire comandi Maven su tutti i progetti selezionati.
Questo comando supporta varie operazioni Maven come install, clean, compile, etc.

Per utilizzare questo comando, è necessario specificare un sottocomando.
Esempi:
  projman mvn install         - Esegue mvn install senza test
  projman mvn install --tests - Esegue mvn install con test abilitati`,
	Run: cmdutil.RequireSubcommandHandler("mvn"),
}

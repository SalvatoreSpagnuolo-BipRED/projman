package git

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/cmdutil"
	"github.com/spf13/cobra"
)

// GitCmd rappresenta il comando parent per tutte le operazioni Git
var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Gestisce le operazioni Git sui progetti selezionati",
	Long: `Esegue operazioni Git su tutti i progetti precedentemente selezionati.
Supporta vari comandi Git come update (pull) per mantenere i progetti aggiornati.

Per utilizzare questo comando, è necessario specificare un sottocomando.
Esempi:
  projman git update    - Aggiorna tutti i progetti con git pull/merge`,
	Run: cmdutil.RequireSubcommandHandler("git"),
}

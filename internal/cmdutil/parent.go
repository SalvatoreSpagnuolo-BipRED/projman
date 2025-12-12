package cmdutil

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// RequireSubcommandHandler restituisce un handler che mostra un errore quando
// un comando parent viene eseguito senza specificare un sottocomando
// Questo elimina la duplicazione presente in cmd/git/git.go e cmd/mvn/mvn.go
func RequireSubcommandHandler(cmdName string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		pterm.Error.Println("Devi specificare un sottocomando")
		pterm.Info.Printf("Usa 'projman %s --help' per vedere i sottocomandi disponibili\n", cmdName)
		_ = cmd.Help()
	}
}

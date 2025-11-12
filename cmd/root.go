package cmd

import (
	"fmt"
	"os"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/cmd/git"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/cmd/mvn"
	"github.com/spf13/cobra"
)

// Version viene impostata durante la build tramite ldflags
var Version = "dev"

// RootCmd rappresenta il comando base quando viene chiamato senza sottocomandi
var RootCmd = &cobra.Command{
	Use:     "projman",
	Short:   "Gestione multipla di progetti Git e Maven",
	Version: Version,
	Long: `Projman è uno strumento da linea di comando che permette di gestire multiple repository Git 
e progetti Maven contemporaneamente. Consente di selezionare un gruppo di progetti 
e eseguire operazioni batch come git pull o mvn install su tutti i progetti selezionati.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Mostra l'help quando viene chiamato senza sottocomandi
		_ = cmd.Help()
	},
}

// Execute esegue il comando root e gestisce l'exit code in caso di errore
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Registra i comandi dei subpackage
	RootCmd.AddCommand(git.GetGitCmd())
	RootCmd.AddCommand(mvn.GetMvnCmd())

	// Personalizza il template della versione
	RootCmd.SetVersionTemplate(fmt.Sprintf("Projman v%s\n", Version))
}

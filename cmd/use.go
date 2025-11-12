package cmd

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// useCmd rappresenta il comando use per cambiare il profilo attivo
var useCmd = &cobra.Command{
	Use:   "use <nome-profilo>",
	Short: "Imposta il profilo corrente",
	Long: `Imposta il profilo specificato come profilo attivo.
Tutti i comandi successivi utilizzeranno la configurazione di questo profilo.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		if err := config.SetCurrentProfile(profileName); err != nil {
			pterm.Error.Println("Errore nell'impostazione del profilo:", err)
			return err
		}

		pterm.Success.Printf("Profilo '%s' impostato come attivo\n", profileName)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(useCmd)
}

package cmd

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// listCmd rappresenta il comando list per visualizzare i profili
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Elenca tutti i profili disponibili",
	Long: `Visualizza la lista di tutti i profili di configurazione salvati.
Indica con un asterisco (*) il profilo attualmente attivo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, current, err := config.ListProfiles()
		if err != nil {
			pterm.Error.Println("Errore nel caricamento dei profili:", err)
			return err
		}

		if len(profiles) == 0 {
			pterm.Info.Println("Nessun profilo configurato")
			pterm.Info.Println("Esegui 'projman init <nome-profilo> <directory>' per creare un nuovo profilo")
			return nil
		}

		pterm.DefaultHeader.Println("Profili disponibili")
		pterm.Println()

		for _, name := range profiles {
			if name == current {
				pterm.Success.Printf("  * %s (attivo)\n", name)
			} else {
				pterm.Info.Printf("    %s\n", name)
			}
		}

		pterm.Println()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

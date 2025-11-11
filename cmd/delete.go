package cmd

import (
	"fmt"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// deleteCmd rappresenta il comando delete per eliminare un profilo
var deleteCmd = &cobra.Command{
	Use:   "delete <nome-profilo>",
	Short: "Elimina un profilo",
	Long: `Elimina un profilo di configurazione esistente.
Se il profilo eliminato era quello attivo, verrà automaticamente selezionato
il primo profilo disponibile.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		// Chiedi conferma prima di eliminare
		confirm := pterm.DefaultInteractiveConfirm
		result, err := confirm.Show(fmt.Sprintf("Sei sicuro di voler eliminare il profilo '%s'?", profileName))
		if err != nil {
			pterm.Error.Println("Errore nella conferma:", err)
			return err
		}

		if !result {
			pterm.Info.Println("Operazione annullata")
			return nil
		}

		if err := config.DeleteProfile(profileName); err != nil {
			pterm.Error.Println("Errore nell'eliminazione del profilo:", err)
			return err
		}

		pterm.Success.Printf("Profilo '%s' eliminato con successo\n", profileName)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

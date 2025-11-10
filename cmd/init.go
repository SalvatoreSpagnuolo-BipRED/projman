package cmd

import (
	"fmt"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/project"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// initCmd rappresenta il comando init per inizializzare la configurazione di projman
var initCmd = &cobra.Command{
	Use:   "init <directory_root_progetti>",
	Short: "Seleziona i progetti da includere nella gestione",
	Long: `Inizializza la configurazione di projman selezionando i progetti da gestire.
Richiede il percorso della directory root contenente i progetti Maven.
Permette di selezionare interattivamente quali progetti includere nella gestione.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Valida che sia stata fornita la directory root
		if len(args) < 1 {
			pterm.Error.Println("Devi specificare la directory root dei progetti")
			pterm.Info.Println("Uso: projman init <directory_root_progetti>")
			return fmt.Errorf("directory root non specificata")
		}

		// Valida l'esistenza e validità della directory
		root, err := config.CheckAndGetDirectory(args[0])
		if err != nil {
			pterm.Error.Println("La directory specificata non è valida:", args[0])
			return err
		}

		// Scansiona la directory alla ricerca di progetti Maven
		pterm.Info.Println("Scansione dei progetti Maven in corso...")
		pterm.Info.Println("Directory:", root)

		projs, err := project.Discover(root)
		if err != nil {
			pterm.Error.Println("Errore durante la scansione dei progetti:", err)
			return err
		}

		// Verifica che siano stati trovati progetti
		if len(projs) == 0 {
			pterm.Warning.Println("Nessun progetto Maven trovato nella directory specificata")
			pterm.Info.Println("Assicurati che la directory contenga sottocartelle con file pom.xml")
			return nil
		}

		pterm.Success.Printf("Trovati %d progetti Maven\n", len(projs))

		// Mostra prompt interattivo per la selezione dei progetti
		names := project.Names(projs)
		selectedNames, err := pterm.DefaultInteractiveMultiselect.
			WithOptions(names).
			WithDefaultOptions(names).
			Show("Seleziona i progetti da includere:")
		if err != nil {
			pterm.Error.Println("Errore durante la selezione dei progetti:", err)
			return err
		}

		// Verifica che almeno un progetto sia stato selezionato
		if len(selectedNames) == 0 {
			pterm.Warning.Println("Nessun progetto selezionato")
			return nil
		}

		// Salva la configurazione nel file JSON
		cfg := config.Config{
			SelectedProjects: selectedNames,
			RootOfProjects:   root,
		}

		if err := config.SaveSettings(cfg); err != nil {
			pterm.Error.Println("Errore durante il salvataggio della configurazione:", err)
			return err
		}

		pterm.Success.Printf("Configurazione completata con %d progetti selezionati\n", len(selectedNames))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

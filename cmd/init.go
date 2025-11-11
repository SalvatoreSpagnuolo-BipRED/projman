package cmd

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/project"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// initCmd rappresenta il comando init per inizializzare la configurazione di projman
var initCmd = &cobra.Command{
	Use:   "init <nome-profilo> <directory_root_progetti>",
	Short: "Seleziona i progetti da includere nella gestione",
	Long: `Inizializza la configurazione di projman selezionando i progetti da gestire.
Richiede il nome del profilo e il percorso della directory root contenente i progetti Maven.
Permette di selezionare interattivamente quali progetti includere nella gestione.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]
		directory := args[1]

		// Valida l'esistenza e validità della directory
		root, err := config.CheckAndGetDirectory(directory)
		if err != nil {
			pterm.Error.Println("La directory specificata non è valida:", directory)
			return err
		}

		// Scansiona la directory alla ricerca di progetti Maven
		pterm.Info.Println("Scansione dei progetti Maven in corso...")
		pterm.Info.Printf("Profilo: %s\n", profileName)
		pterm.Info.Printf("Directory: %s\n", root)

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

		// Salva la configurazione nel profilo specificato
		cfg := config.Config{
			SelectedProjects: selectedNames,
			RootOfProjects:   root,
		}

		if err := config.SaveProfile(profileName, cfg); err != nil {
			pterm.Error.Println("Errore durante il salvataggio della configurazione:", err)
			return err
		}

		pterm.Success.Printf("Profilo '%s' configurato con %d progetti selezionati\n", profileName, len(selectedNames))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

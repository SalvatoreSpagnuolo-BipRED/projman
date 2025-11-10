package cmd

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/project"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Seleziona i progetti da includere nella gestione",
	Long: `Inizializza la configurazione di projman selezionando i progetti da gestire.
Richiede il percorso della directory root contenente i progetti Maven.
Permette di selezionare interattivamente quali progetti includere nella gestione.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Controlla che sia stata passata la directory root
		if len(args) < 1 {
			pterm.Error.Println("Devi specificare la directory root dei progetti")
			return nil
		}

		// Controlla che la directory esista
		root, err := config.CheckAndGetDirectory(args[0])
		if err != nil {
			pterm.Error.Println("Immetti una directory valida")
			return nil
		}

		// Scansione dei progetti Maven
		pterm.Info.Println("Scansione dei progetti in " + root)
		projs, _ := project.Discover(root)
		if len(projs) == 0 {
			pterm.Warning.Println("Nessun progetto Maven trovato")
			return nil
		}

		// Prompt per la selezione dei progetti
		names := project.Names(projs)
		selectedNames, err := pterm.DefaultInteractiveMultiselect.
			WithOptions(names).
			WithDefaultOptions(names).
			Show("Seleziona i progetti da includere:")
		if err != nil {
			return err
		}

		// Salva lo config in json
		cfg := config.Config{
			SelectedProjects: selectedNames,
			RootOfProjects:   root,
		}
		return config.SaveSettings(cfg)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

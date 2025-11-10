package mvn

import (
	"fmt"
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var runTests bool

// installCmd rappresenta il comando per eseguire mvn install sui progetti selezionati
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Esegue mvn install sui progetti selezionati",
	Long: `Esegue il comando 'mvn install' su tutti i progetti Maven selezionati.
Per default i test sono disabilitati. Usa il flag --tests o -t per abilitarli.
Il comando cerca il file pom.xml in ogni progetto selezionato ed esegue l'installazione.

Esempi:
  projman mvn install         - Installa i progetti senza eseguire i test
  projman mvn install --tests - Installa i progetti eseguendo i test`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Maven Install")

		// Carica la configurazione salvata
		cfg, err := config.LoadSettings()
		if err != nil {
			pterm.Error.Println("Errore nel caricamento della configurazione:", err)
			pterm.Info.Println("Esegui prima 'projman init <directory>' per configurare i progetti")
			return
		}

		// Verifica che ci siano progetti selezionati
		if len(cfg.SelectedProjects) == 0 {
			pterm.Warning.Println("Nessun progetto selezionato nella configurazione")
			pterm.Info.Println("Esegui 'projman init <directory>' per selezionare i progetti")
			return
		}

		// Mostra informazioni sull'esecuzione
		if runTests {
			pterm.Info.Println("Esecuzione con test abilitati")
		} else {
			pterm.Info.Println("Esecuzione con test disabilitati (-DskipTests=true)")
		}

		// Contatori per statistiche finali
		successCount := 0
		failureCount := 0

		// Esegui mvn install per ogni progetto selezionato
		spinner, _ := pterm.DefaultSpinner.Start("Inizializzazione...")

		for _, projectName := range cfg.SelectedProjects {
			spinner.UpdateText(fmt.Sprintf("Installazione progetto: %s", projectName))

			// Costruisci il percorso del pom.xml
			pomPath := filepath.Join(cfg.RootOfProjects, projectName, "pom.xml")

			// Prepara gli argomenti per Maven
			args := buildMavenArgs(pomPath, runTests)

			// Esegui il comando Maven
			if err := exec.Run("mvn", args...); err != nil {
				pterm.Error.Printf("✗ Progetto %s: errore durante l'installazione\n", projectName)
				pterm.Error.Println("  ", err.Error())
				failureCount++
			} else {
				pterm.Success.Printf("✓ Progetto %s: installazione completata\n", projectName)
				successCount++
			}
		}

		_ = spinner.Stop()

		// Mostra il riepilogo finale
		printInstallSummary(successCount, failureCount)
	},
}

// buildMavenArgs costruisce gli argomenti per il comando Maven
func buildMavenArgs(pomPath string, includeTests bool) []string {
	args := []string{"-f", pomPath, "install"}
	if !includeTests {
		args = append(args, "-DskipTests=true")
	}
	return args
}

// printInstallSummary stampa un riepilogo delle operazioni di installazione
func printInstallSummary(successCount, failureCount int) {
	pterm.Println()
	pterm.DefaultSection.Println("Riepilogo Operazioni")

	totalCount := successCount + failureCount
	pterm.Info.Printf("Progetti totali: %d\n", totalCount)

	if successCount > 0 {
		pterm.Success.Printf("Installazioni riuscite: %d\n", successCount)
	}

	if failureCount > 0 {
		pterm.Error.Printf("Installazioni fallite: %d\n", failureCount)
	} else {
		pterm.Success.Println("Tutte le installazioni completate con successo!")
	}
}

func init() {
	mvnCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&runTests, "tests", "t", false, "Abilita l'esecuzione dei test durante l'installazione")
}

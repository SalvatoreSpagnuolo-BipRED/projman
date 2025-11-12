package mvn

import (
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/maven"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/maven/executor"
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

		// Carica e valida la configurazione
		cfg, err := config.LoadAndValidateConfig()
		if err != nil {
			return
		}

		// Permette all'utente di selezionare i progetti da aggiornare
		selectedProjects, err := config.SelectProjectsToUpdate(cfg)
		if err != nil {
			return
		}

		// Salva la selezione aggiornata
		cfg.SelectedProjects = selectedProjects
		if err := config.SaveSettings(*cfg); err != nil {
			pterm.Error.Println("Errore nel salvataggio della configurazione:", err)
			return
		}

		// Mostra informazioni sull'esecuzione
		if runTests {
			pterm.Info.Println("Esecuzione con test abilitati")
		} else {
			pterm.Info.Println("Esecuzione con test disabilitati (-DskipTests=true)")
		}

		// Costruisci il grafo delle dipendenze
		spinner, _ := pterm.DefaultSpinner.Start("Analisi dipendenze Maven...")

		dependencyGraph, err := maven.BuildDependencyGraph(cfg.SelectedProjects, cfg.RootOfProjects)
		if err != nil {
			spinner.Fail("Errore durante l'analisi delle dipendenze:", err)
			return
		}

		// Ordina i progetti topologicamente in base alle dipendenze
		sortedProjects, err := maven.TopologicalSort(dependencyGraph)
		if err != nil {
			spinner.Fail("Errore durante l'ordinamento dei progetti:", err)
			pterm.Error.Println("Le dipendenze tra i progetti potrebbero contenere cicli")
			return
		}

		spinner.Success("Analisi dipendenze completata")

		// Mostra l'ordine di installazione
		pterm.Info.Println("\nOrdine di installazione (basato sulle dipendenze):")
		for i, projectName := range sortedProjects {
			deps := dependencyGraph[projectName]
			if len(deps) > 0 {
				pterm.Info.Printf("  %d. %s (dipende da: %v)\n", i+1, projectName, deps)
			} else {
				pterm.Info.Printf("  %d. %s (nessuna dipendenza interna)\n", i+1, projectName)
			}
		}
		pterm.Println()

		// Contatori per statistiche finali
		successCount := 0
		failureCount := 0
		skippedCount := 0

		// Esegui mvn install per ogni progetto nell'ordine corretto
		for i, projectName := range sortedProjects {
			// Mostra un'intestazione per il progetto corrente
			pterm.DefaultHeader.WithFullWidth().Printf("Progetto %d/%d: %s", i+1, len(sortedProjects), projectName)

			// Costruisci il percorso del pom.xml
			pomPath := filepath.Join(cfg.RootOfProjects, projectName, "pom.xml")

			// Prepara gli argomenti per Maven
			args := buildMavenArgs(pomPath, runTests)

			// Esegui il comando Maven con il nuovo executor
			mavenExec := executor.NewMavenExecutor(projectName, args)
			if err := mavenExec.Run(); err != nil {
				pterm.Error.Printf("✗ Progetto %s: errore durante l'installazione\n", projectName)
				pterm.Error.Println("  ", err.Error())
				failureCount++

				// Chiedi all'utente se vuole continuare
				if !exec.WaitForUserInput(projectName) {
					pterm.Warning.Println("Esecuzione interrotta dall'utente")
					skippedCount = len(sortedProjects) - i - 1
					break
				}
			} else {
				successCount++
			}

			// Aggiungi una riga vuota tra i progetti (tranne dopo l'ultimo)
			if i < len(sortedProjects)-1 {
				pterm.Println()
			}
		}

		// Mostra il riepilogo finale
		printInstallSummary(successCount, failureCount, skippedCount)
	},
}

// buildMavenArgs costruisce gli argomenti per il comando Maven
func buildMavenArgs(pomPath string, includeTests bool) []string {
	args := []string{"-B", "-f", pomPath, "install"}
	if !includeTests {
		args = append(args, "-DskipTests=true")
	}
	return args
}

// printInstallSummary stampa un riepilogo delle operazioni di installazione
func printInstallSummary(successCount, failureCount, skippedCount int) {
	pterm.Println()
	pterm.DefaultSection.Println("Riepilogo Operazioni")

	totalCount := successCount + failureCount + skippedCount
	pterm.Info.Printf("Progetti totali: %d\n", totalCount)

	if successCount > 0 {
		pterm.Success.Printf("Installazioni riuscite: %d\n", successCount)
	}

	if failureCount > 0 {
		pterm.Error.Printf("Installazioni fallite: %d\n", failureCount)
	}

	if skippedCount > 0 {
		pterm.Warning.Printf("Installazioni saltate: %d\n", skippedCount)
	}

	if failureCount == 0 && skippedCount == 0 {
		pterm.Success.Println("Tutte le installazioni completate con successo!")
	}
}

func init() {
	mvnCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&runTests, "tests", "t", false, "Abilita l'esecuzione dei test durante l'installazione")
}

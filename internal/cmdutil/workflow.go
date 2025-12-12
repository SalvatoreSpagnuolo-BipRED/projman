// Package cmdutil fornisce utilities comuni per i comandi CLI
package cmdutil

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec"
	"github.com/pterm/pterm"
)

// LoadConfigAndSelectProjects combina il pattern ripetuto di:
// 1. Caricamento e validazione configurazione
// 2. Selezione interattiva dei progetti
// 3. Salvataggio della selezione
// Questo elimina la duplicazione presente in git/update.go e mvn/install.go
func LoadConfigAndSelectProjects() (*config.Config, []string, error) {
	// Carica e valida la configurazione
	cfg, err := config.LoadSettings()
	if err != nil {
		pterm.Error.Println("Errore nel caricamento della configurazione:", err)
		pterm.Info.Println("Esegui prima 'projman init <nome-profilo> <directory>' per inizializzare")
		return nil, nil, err
	}

	// Permette all'utente di selezionare i progetti da aggiornare
	selectedProjects, err := config.SelectProjectsToUpdate(&cfg)
	if err != nil {
		return nil, nil, err
	}

	// Salva la selezione aggiornata
	cfg.SelectedProjects = selectedProjects
	if err := config.SaveSettings(cfg); err != nil {
		pterm.Error.Println("Errore nel salvataggio della configurazione:", err)
		return nil, nil, err
	}

	return &cfg, selectedProjects, nil
}

// ProjectProcessor è una funzione che processa un singolo progetto
// Riceve il nome del progetto, l'indice corrente e il numero totale di progetti
type ProjectProcessor func(projectName string, index int, total int) error

// ProcessProjectsWithErrorHandling esegue il processor per ogni progetto con gestione errori unificata
// Se un progetto fallisce, chiede all'utente se vuole continuare o interrompere
// Restituisce il numero di progetti processati con successo, falliti e saltati
func ProcessProjectsWithErrorHandling(projects []string, processor ProjectProcessor) (success, failed, skipped int) {
	for i, projectName := range projects {
		if err := processor(projectName, i, len(projects)); err != nil {
			pterm.Error.Printf("✗ Errore durante l'elaborazione di '%s'\n", projectName)
			failed++

			if !exec.WaitForUserInput(projectName) {
				pterm.Warning.Println("Esecuzione interrotta dall'utente")
				skipped = len(projects) - i - 1
				return
			}
		} else {
			success++
		}

		// Aggiungi riga vuota tra progetti (eccetto l'ultimo)
		if i < len(projects)-1 {
			pterm.Println()
		}
	}

	return
}

// OperationSummary contiene le statistiche delle operazioni eseguite
type OperationSummary struct {
	Total   int
	Success int
	Failed  int
	Skipped int
}

// PrintOperationSummary stampa un riepilogo delle operazioni eseguite
func PrintOperationSummary(summary OperationSummary) {
	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println("RIEPILOGO")

	pterm.Info.Printf("Totale progetti: %d\n", summary.Total)

	if summary.Success > 0 {
		pterm.Success.Printf("✓ Successo: %d\n", summary.Success)
	}

	if summary.Failed > 0 {
		pterm.Error.Printf("✗ Falliti: %d\n", summary.Failed)
	}

	if summary.Skipped > 0 {
		pterm.Warning.Printf("⊘ Saltati: %d\n", summary.Skipped)
	}
}

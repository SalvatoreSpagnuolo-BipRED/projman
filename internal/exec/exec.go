package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec/scrollable"
)

// Run esegue un comando esterno mostrando l'output in tempo reale.
// Stampa il comando prima di eseguirlo e collega stdout/stderr al terminale corrente.
// Restituisce un errore se il comando fallisce.
func Run(name string, args ...string) error {
	// Mostra il comando che sta per essere eseguito
	fmt.Printf("$ %s %s\n", name, strings.Join(args, " "))

	// Crea il comando con gli argomenti forniti
	cmd := exec.Command(name, args...)

	// Collega l'output del comando al terminale corrente
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Passa le variabili d'ambiente correnti al comando
	cmd.Env = os.Environ()

	// Esegue il comando e restituisce l'eventuale errore
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("comando fallito '%s %s': %w", name, strings.Join(args, " "), err)
	}

	return nil
}

// RunWithOutput esegue un comando esterno e restituisce l'output come stringa.
// A differenza di Run, non mostra l'output in tempo reale ma lo cattura per restituirlo.
// L'output viene trimmato degli spazi bianchi iniziali e finali.
// Restituisce l'output e un eventuale errore.
func RunWithOutput(name string, args ...string) (string, error) {
	// Crea il comando con gli argomenti forniti
	cmd := exec.Command(name, args...)

	// Passa le variabili d'ambiente correnti al comando
	cmd.Env = os.Environ()

	// Esegue il comando e cattura l'output
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("comando fallito '%s %s': %w", name, strings.Join(args, " "), err)
	}

	// Restituisce l'output come stringa, rimuovendo spazi bianchi
	return strings.TrimSpace(string(output)), nil
}

// RunWithScrollableOutput esegue un comando esterno mostrando l'output in un'area scrollabile.
// L'area ha un'altezza fissa definita da maxLines, e l'output scorre verso il basso
// man mano che vengono ricevute nuove righe, con un indicatore di progresso animato.
func RunWithScrollableOutput(name string, args []string, maxLines int) error {
	// Setup iniziale del comando
	executor := scrollable.NewExecutor(name, args)
	executor.PrintCommandInfo()
	executor.PrintSectionHeader(maxLines)

	// Prepara il comando con le pipe
	if err := executor.SetupCommand(); err != nil {
		return err
	}

	// Ottieni il comando preparato
	cmd := executor.GetCommand()

	// Ottieni le pipe per stdout e stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("errore creazione stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("errore creazione stderr pipe: %w", err)
	}

	// Passa le variabili d'ambiente
	cmd.Env = os.Environ()

	// Avvia il comando
	if err := executor.Start(); err != nil {
		return err
	}

	// Inizializza i componenti per la gestione dell'output
	startTime := time.Now()
	area, err := scrollable.NewArea(maxLines)
	if err != nil {
		return err
	}
	defer area.Stop()

	buffer := scrollable.NewBuffer(maxLines)

	// Funzione di aggiornamento dell'area
	updateArea := func() {
		area.Update(buffer.GetLines())
	}

	// Avvia l'aggiornamento periodico
	updater := scrollable.NewUpdater(500*time.Millisecond, updateArea)
	updater.Start()
	defer updater.Stop()

	// Gestisci lo streaming dell'output
	streamer := scrollable.NewStreamer(buffer, updateArea)
	go streamer.StreamOutput(stdout)
	go streamer.StreamOutput(stderr)
	streamer.WaitForCompletion()

	// Attendi il completamento del comando
	cmdErr := executor.Wait()

	// Stampa il messaggio di completamento se successo
	if cmdErr == nil {
		executor.PrintCompletionMessage(startTime)
	}

	return cmdErr
}

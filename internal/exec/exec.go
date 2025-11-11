package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pterm/pterm"
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

// RunWithSpinner esegue un comando esterno mostrando uno spinner semplice con timer.
// Lo spinner mostra solo il tempo trascorso, senza log di output (evita artefatti su terminali lenti).
// Il parametro maxLogLines viene ignorato ed è mantenuto solo per retro-compatibilità.
func RunWithSpinner(name string, args []string, _ int) error {
	// Mostra il comando che sta per essere eseguito
	pterm.Info.Printf("$ %s %s\n", name, strings.Join(args, " "))

	// Prepara il comando
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()

	// Redireziona stdout e stderr a /dev/null (o equivalente) per evitare output
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Avvia il comando
	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		pterm.Error.Println("Errore avvio comando")
		return fmt.Errorf("errore avvio comando '%s %s': %w", name, strings.Join(args, " "), err)
	}

	// Crea lo spinner che girerà continuamente con messaggio fisso
	spinnerText := "Esecuzione in corso"
	spinner, _ := pterm.DefaultSpinner.Start(spinnerText)

	// Area per il timer
	area, _ := pterm.DefaultArea.Start()

	// Goroutine per aggiornare il timer ogni secondo
	stopUpdater := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)
				minutes := int(elapsed.Minutes())
				seconds := int(elapsed.Seconds()) % 60
				area.Update(fmt.Sprintf("Tempo trascorso: %dm %ds", minutes, seconds))
			case <-stopUpdater:
				return
			}
		}
	}()

	// Attendi il completamento del comando
	cmdErr := cmd.Wait()

	// Ferma l'updater e l'area
	close(stopUpdater)
	_ = area.Stop()

	// Mostra il risultato
	elapsed := time.Since(startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60

	if cmdErr != nil {
		spinner.Fail(fmt.Sprintf("Comando fallito dopo %dm %ds", minutes, seconds))
		return fmt.Errorf("comando fallito '%s %s': %w", name, strings.Join(args, " "), cmdErr)
	}

	spinner.Success(fmt.Sprintf("Comando completato in %dm %ds", minutes, seconds))
	return nil
}

// RunWithScrollableOutput è un alias per RunWithSpinner per retro-compatibilità.
// Il parametro maxLines specifica quante righe di log mostrare (usa 0 per nessun log).
func RunWithScrollableOutput(name string, args []string, maxLines int) error {
	return RunWithSpinner(name, args, maxLines)
}

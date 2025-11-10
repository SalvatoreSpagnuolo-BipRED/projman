package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

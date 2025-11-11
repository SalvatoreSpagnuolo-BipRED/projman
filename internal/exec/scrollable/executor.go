package scrollable

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// Executor gestisce l'esecuzione di un comando esterno
type Executor struct {
	name string
	args []string
	cmd  *exec.Cmd
}

// NewExecutor crea un nuovo executor per il comando specificato
func NewExecutor(name string, args []string) *Executor {
	return &Executor{
		name: name,
		args: args,
	}
}

// PrintCommandInfo stampa le informazioni sul comando prima dell'esecuzione
func (ce *Executor) PrintCommandInfo() {
	pterm.Info.Printf("$ %s %s\n", ce.name, strings.Join(ce.args, " "))
}

// PrintSectionHeader stampa l'intestazione della sezione di output
func (ce *Executor) PrintSectionHeader(maxLines int) {
	pterm.DefaultSection.Printf("Output comando (ultime %d righe)", maxLines)
}

// SetupCommand prepara il comando con le pipe per stdout e stderr
func (ce *Executor) SetupCommand() error {
	ce.cmd = exec.Command(ce.name, ce.args...)
	return nil
}

// GetCommand restituisce il comando preparato
func (ce *Executor) GetCommand() *exec.Cmd {
	return ce.cmd
}

// Start avvia il comando
func (ce *Executor) Start() error {
	if err := ce.cmd.Start(); err != nil {
		return fmt.Errorf("errore avvio comando '%s %s': %w", ce.name, strings.Join(ce.args, " "), err)
	}
	return nil
}

// Wait attende il completamento del comando
func (ce *Executor) Wait() error {
	if err := ce.cmd.Wait(); err != nil {
		return fmt.Errorf("comando fallito '%s %s': %w", ce.name, strings.Join(ce.args, " "), err)
	}
	return nil
}

// PrintCompletionMessage stampa il messaggio di completamento con il tempo trascorso
func (ce *Executor) PrintCompletionMessage(startTime time.Time) {
	elapsed := time.Since(startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	pterm.Print("\r" + strings.Repeat(" ", 80) + "\r")
	pterm.Success.Printf("Comando completato in %dm %ds\n", minutes, seconds)
	pterm.Println()
}

package scrollable

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// Area gestisce un'area scrollabile per l'output dei comandi
type Area struct {
	area      *pterm.AreaPrinter
	maxLines  int
	startTime time.Time
}

// NewArea crea e avvia una nuova area scrollabile
func NewArea(maxLines int) (*Area, error) {
	area, err := pterm.DefaultArea.Start()
	if err != nil {
		return nil, fmt.Errorf("errore creazione area: %w", err)
	}

	return &Area{
		area:      area,
		maxLines:  maxLines,
		startTime: time.Now(),
	}, nil
}

// Update aggiorna il contenuto dell'area con le righe fornite
func (sa *Area) Update(lines []string) {
	displayLines := sa.buildDisplayLines(lines)
	sa.area.Update(strings.Join(displayLines, "\n"))
}

// buildDisplayLines costruisce le righe da visualizzare con spinner e tempo
func (sa *Area) buildDisplayLines(lines []string) []string {
	displayLines := make([]string, 0, sa.maxLines+2)

	// Aggiungi le righe di output (troncate se necessario)
	for _, line := range lines {
		displayLines = append(displayLines, sa.truncateLine(line))
	}

	// Riempi con righe vuote per mantenere l'altezza fissa
	for len(displayLines) < sa.maxLines {
		displayLines = append(displayLines, "")
	}

	// Aggiungi separatore e indicatore di progresso
	displayLines = append(displayLines, sa.buildSeparator())
	displayLines = append(displayLines, sa.buildStatusLine())

	return displayLines
}

// truncateLine tronca una riga se supera la lunghezza massima
func (sa *Area) truncateLine(line string) string {
	const maxLineLength = 120
	if len(line) > maxLineLength {
		return line[:maxLineLength-3] + "..."
	}
	return line
}

// buildSeparator costruisce la riga separatrice
func (sa *Area) buildSeparator() string {
	return strings.Repeat("─", 60)
}

// buildStatusLine costruisce la riga di stato con spinner e tempo
func (sa *Area) buildStatusLine() string {
	spinChar := sa.getSpinnerChar()
	minutes, seconds := sa.getElapsedTime()
	return fmt.Sprintf("%s Esecuzione in corso... (tempo: %dm %ds)", spinChar, minutes, seconds)
}

// getSpinnerChar restituisce il carattere dello spinner corrente
func (sa *Area) getSpinnerChar() string {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinIndex := int(time.Since(sa.startTime).Milliseconds()/100) % len(spinChars)
	return spinChars[spinIndex]
}

// getElapsedTime calcola e restituisce il tempo trascorso in minuti e secondi
func (sa *Area) getElapsedTime() (minutes, seconds int) {
	elapsed := time.Since(sa.startTime)
	minutes = int(elapsed.Minutes())
	seconds = int(elapsed.Seconds()) % 60
	return
}

// Stop ferma l'area e pulisce la riga corrente
func (sa *Area) Stop() {
	_ = sa.area.Stop()
	// Pulisci la riga corrente completamente
	fmt.Print("\r" + strings.Repeat(" ", 100) + "\r")
}

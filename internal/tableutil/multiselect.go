package tableutil

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// MultiSelectTable rappresenta una tabella interattiva con multiselect
type MultiSelectTable struct {
	Headers        []string
	Rows           [][]string
	Message        string
	DefaultIndices []int
}

// Show mostra la tabella multiselect e restituisce gli indici delle righe selezionate
func (t *MultiSelectTable) Show() ([]int, error) {
	if len(t.Rows) == 0 {
		return nil, fmt.Errorf("nessuna riga da visualizzare")
	}

	if len(t.Headers) == 0 {
		return nil, fmt.Errorf("nessun header specificato")
	}

	// Calcola le larghezze massime per ogni colonna
	colWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		colWidths[i] = len(header)
	}

	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Costruisci l'header formattato
	headerParts := make([]string, len(t.Headers))
	for i, header := range t.Headers {
		headerParts[i] = fmt.Sprintf("%-*s", colWidths[i], header)
	}
	headerLine := strings.Repeat(" ", 6) + strings.Join(headerParts, " │ ")
	separatorLine := strings.Repeat("─", len(headerLine))

	pterm.Info.Println(t.Message)
	pterm.Println()

	// Crea il messaggio con l'header della tabella
	message := fmt.Sprintf("%s\n%s", headerLine, separatorLine)

	// Formatta le righe come opzioni per la multiselect in formato tabellare
	options := make([]string, len(t.Rows))
	for i, row := range t.Rows {
		rowParts := make([]string, len(t.Headers))
		for j, cell := range row {
			if j < len(colWidths) {
				rowParts[j] = fmt.Sprintf("%-*s", colWidths[j], cell)
			}
		}
		options[i] = strings.Join(rowParts, " │ ")
	}

	// Converti gli indici di default alle opzioni corrispondenti
	defaultOptions := make([]string, 0)
	for _, idx := range t.DefaultIndices {
		if idx >= 0 && idx < len(options) {
			defaultOptions = append(defaultOptions, options[idx])
		}
	}

	// Multiselect interattiva con formato tabellare
	selectedOptions, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithDefaultOptions(defaultOptions).
		Show(message)
	if err != nil {
		return nil, err
	}

	// Converti le opzioni selezionate negli indici
	selectedIndices := make([]int, 0)
	for _, selectedOption := range selectedOptions {
		for i, option := range options {
			if option == selectedOption {
				selectedIndices = append(selectedIndices, i)
				break
			}
		}
	}

	return selectedIndices, nil
}

package tableutil

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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

	// Costruisci l'header della tabella
	headerParts := make([]string, len(t.Headers))
	for i, header := range t.Headers {
		headerParts[i] = fmt.Sprintf("%-*s", colWidths[i], header)
	}
	headerLine := "       " + strings.Join(headerParts, "  │  ")

	// Formatta le righe come opzioni per la multiselect
	options := make([]string, len(t.Rows))
	for i, row := range t.Rows {
		rowParts := make([]string, len(t.Headers))
		for j, cell := range row {
			if j < len(colWidths) {
				rowParts[j] = fmt.Sprintf("%-*s", colWidths[j], cell)
			}
		}
		options[i] = strings.Join(rowParts, "  │  ")
	}

	// Messaggio di default
	message := fmt.Sprintf("%s\n%s\n%s", t.Message, headerLine, strings.Repeat("─", len(headerLine)))
	if message == "" {
		message = "Usa ↑↓ per navigare, Spazio per selezionare, Invio per confermare:"
	}

	// Multiselect interattiva
	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Default: t.DefaultIndices,
		Help:    "",
	}

	err := survey.AskOne(prompt, &selectedIndices)
	if err != nil {
		return nil, err
	}

	return selectedIndices, nil
}

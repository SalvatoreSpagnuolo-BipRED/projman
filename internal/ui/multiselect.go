package ui

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

const (
	// TableCellPadding è il padding usato per l'indentazione della tabella
	TableCellPadding = 6
	// TableColumnSeparator è il carattere usato per separare le colonne
	TableColumnSeparator = " │ "
	// TableHeaderSeparatorChar è il carattere usato per la linea di separazione dell'header
	TableHeaderSeparatorChar = "─"
)

// MultiSelectTable rappresenta una tabella interattiva con selezione multipla.
// Permette all'utente di selezionare più righe da una tabella formattata.
type MultiSelectTable struct {
	Headers        []string   // Intestazioni delle colonne
	Rows           [][]string // Righe della tabella (ogni riga è un array di celle)
	Message        string     // Messaggio da mostrare prima della tabella
	DefaultIndices []int      // Indici delle righe selezionate di default
}

// Show mostra la tabella interattiva e restituisce gli indici delle righe selezionate dall'utente
func (t *MultiSelectTable) Show() ([]int, error) {
	// Valida i dati della tabella
	if err := t.validate(); err != nil {
		return nil, err
	}

	// Calcola le larghezze delle colonne
	colWidths := t.calculateColumnWidths()

	// Costruisce l'header formattato
	headerLine := t.buildHeaderLine(colWidths)
	separatorLine := strings.Repeat(TableHeaderSeparatorChar, len(headerLine))

	// Mostra il messaggio e l'header
	pterm.Info.Println(t.Message)
	pterm.Println()

	message := fmt.Sprintf("%s\n%s", headerLine, separatorLine)

	// Formatta le righe come opzioni per la selezione multipla
	options := t.buildRowOptions(colWidths)

	// Converti gli indici di default alle opzioni corrispondenti
	defaultOptions := t.getDefaultOptions(options)

	// Mostra la selezione interattiva multipla
	selectedOptions, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithDefaultOptions(defaultOptions).
		Show(message)
	if err != nil {
		return nil, fmt.Errorf("errore durante la selezione interattiva: %w", err)
	}

	// Converti le opzioni selezionate negli indici corrispondenti
	return t.optionsToIndices(selectedOptions, options), nil
}

// validate verifica che i dati della tabella siano validi
func (t *MultiSelectTable) validate() error {
	if len(t.Rows) == 0 {
		return fmt.Errorf("nessuna riga da visualizzare")
	}

	if len(t.Headers) == 0 {
		return fmt.Errorf("nessun header specificato")
	}

	return nil
}

// calculateColumnWidths calcola la larghezza massima necessaria per ogni colonna
func (t *MultiSelectTable) calculateColumnWidths() []int {
	colWidths := make([]int, len(t.Headers))

	// Inizializza con le larghezze degli header
	for i, header := range t.Headers {
		colWidths[i] = len(header)
	}

	// Aggiorna con le larghezze massime delle celle
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	return colWidths
}

// buildHeaderLine costruisce la riga di intestazione formattata
func (t *MultiSelectTable) buildHeaderLine(colWidths []int) string {
	headerParts := make([]string, len(t.Headers))
	for i, header := range t.Headers {
		headerParts[i] = fmt.Sprintf("%-*s", colWidths[i], header)
	}

	padding := strings.Repeat(" ", TableCellPadding)
	return padding + strings.Join(headerParts, TableColumnSeparator)
}

// buildRowOptions formatta tutte le righe come opzioni per la selezione multipla
func (t *MultiSelectTable) buildRowOptions(colWidths []int) []string {
	options := make([]string, len(t.Rows))

	for i, row := range t.Rows {
		options[i] = t.formatRow(row, colWidths)
	}

	return options
}

// formatRow formatta una singola riga in formato tabellare
func (t *MultiSelectTable) formatRow(row []string, colWidths []int) string {
	rowParts := make([]string, len(t.Headers))

	for j, cell := range row {
		if j < len(colWidths) {
			rowParts[j] = fmt.Sprintf("%-*s", colWidths[j], cell)
		}
	}

	return strings.Join(rowParts, TableColumnSeparator)
}

// getDefaultOptions converte gli indici di default nelle opzioni corrispondenti
func (t *MultiSelectTable) getDefaultOptions(options []string) []string {
	defaultOptions := make([]string, 0, len(t.DefaultIndices))

	for _, idx := range t.DefaultIndices {
		if idx >= 0 && idx < len(options) {
			defaultOptions = append(defaultOptions, options[idx])
		}
	}

	return defaultOptions
}

// optionsToIndices converte le opzioni selezionate negli indici corrispondenti
func (t *MultiSelectTable) optionsToIndices(selectedOptions, allOptions []string) []int {
	selectedIndices := make([]int, 0, len(selectedOptions))

	for _, selectedOption := range selectedOptions {
		for i, option := range allOptions {
			if option == selectedOption {
				selectedIndices = append(selectedIndices, i)
				break
			}
		}
	}

	return selectedIndices
}

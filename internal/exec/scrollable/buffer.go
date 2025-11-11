package scrollable

import "sync"

// Buffer è un buffer circolare thread-safe per mantenere le ultime N righe
type Buffer struct {
	lines    []string
	maxLines int
	mutex    sync.Mutex
}

// NewBuffer crea un nuovo buffer con capacità massima specificata
func NewBuffer(maxLines int) *Buffer {
	return &Buffer{
		lines:    make([]string, 0, maxLines),
		maxLines: maxLines,
	}
}

// Append aggiunge una nuova riga al buffer, rimuovendo la più vecchia se necessario
func (ob *Buffer) Append(line string) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	ob.lines = append(ob.lines, line)

	// Mantieni solo le ultime maxLines righe
	if len(ob.lines) > ob.maxLines {
		ob.lines = ob.lines[1:]
	}
}

// GetLines restituisce una copia delle righe correnti
func (ob *Buffer) GetLines() []string {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	// Crea una copia per evitare race condition
	linesCopy := make([]string, len(ob.lines))
	copy(linesCopy, ob.lines)
	return linesCopy
}

// Clear svuota il buffer
func (ob *Buffer) Clear() {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	ob.lines = make([]string, 0, ob.maxLines)
}

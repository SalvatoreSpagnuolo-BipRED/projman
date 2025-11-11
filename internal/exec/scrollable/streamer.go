package scrollable

import (
	"bufio"
	"io"
)

// Streamer gestisce la lettura in streaming di stdout e stderr
type Streamer struct {
	buffer     *Buffer
	updateFunc func()
	done       chan bool
}

// NewStreamer crea un nuovo streamer di output
func NewStreamer(buffer *Buffer, updateFunc func()) *Streamer {
	return &Streamer{
		buffer:     buffer,
		updateFunc: updateFunc,
		done:       make(chan bool, 2), // 2 goroutine: stdout e stderr
	}
}

// StreamOutput legge l'output da un reader e lo aggiunge al buffer
func (os *Streamer) StreamOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		os.buffer.Append(line)
		os.updateFunc()
	}

	os.done <- true
}

// WaitForCompletion attende che entrambi i stream (stdout e stderr) siano completati
func (os *Streamer) WaitForCompletion() {
	<-os.done
	<-os.done
}

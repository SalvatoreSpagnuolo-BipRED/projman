package scrollable

import "time"

// Updater gestisce l'aggiornamento periodico dell'area
type Updater struct {
	updateFunc func()
	ticker     *time.Ticker
	stop       chan bool
}

// NewUpdater crea un nuovo updater con intervallo specificato
func NewUpdater(interval time.Duration, updateFunc func()) *Updater {
	return &Updater{
		updateFunc: updateFunc,
		ticker:     time.NewTicker(interval),
		stop:       make(chan bool),
	}
}

// Start avvia l'aggiornamento periodico in una goroutine
func (pu *Updater) Start() {
	go func() {
		defer pu.ticker.Stop()

		for {
			select {
			case <-pu.ticker.C:
				pu.updateFunc()
			case <-pu.stop:
				return
			}
		}
	}()
}

// Stop ferma l'aggiornamento periodico
func (pu *Updater) Stop() {
	close(pu.stop)
}

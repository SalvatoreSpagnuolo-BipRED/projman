/*
Package main è l'entry point dell'applicazione Projman.
Projman è un tool CLI per la gestione batch di progetti Maven/Git che consente
di eseguire operazioni su multipli repository contemporaneamente.

Credits 2025 SALVATORE SPAGNUOLO
*/
package main

import "github.com/SalvatoreSpagnuolo-BipRED/projman/cmd"

// main è il punto di ingresso dell'applicazione.
// Delega l'esecuzione al package cmd che gestisce i comandi CLI tramite Cobra.
func main() {
	cmd.Execute()
}

/*
Copyright © 2025 Salvatore Spagnuolo
*/
package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help [comando]",
	Short: "Mostra l'help di projman o di un comando specifico",
	Long: `Mostra informazioni sull'uso dei comandi di projman.
Se viene passato il nome di un comando, mostra informazioni specifiche su quel comando.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Titolo principale
		pterm.DefaultHeader.WithFullWidth().Println("Projman - Gestore Multiplo di Progetti Git e Maven")
		pterm.Println()

		// Descrizione
		pterm.DefaultBox.WithTitle("DESCRIZIONE").WithTitleTopCenter().WithPadding(1).Println(
			"Projman è un tool CLI per gestire batch di progetti Maven/Git contemporaneamente.\n" +
				"Permette di eseguire operazioni Git e Maven su multipli repository in modo efficiente.",
		)
		pterm.Println()

		// Uso generale
		pterm.DefaultSection.Println("USO")
		pterm.FgLightCyan.Println("  projman [comando] [sottocomando] [argomenti] [flags]")
		pterm.Println()

		// Tabella comandi principali
		pterm.DefaultSection.Println("COMANDI DISPONIBILI")
		tableData := pterm.TableData{
			{"COMANDO", "DESCRIZIONE"},
			{"init [directory]", "Scansiona la directory e seleziona i progetti Maven da gestire"},
			{"git", "Gestisce le operazioni Git sui progetti selezionati"},
			{"git update", "Esegue git pull/merge su tutti i progetti con gestione intelligente dei branch"},
			{"mvn", "Esegue comandi Maven sui progetti selezionati"},
			{"mvn install", "Esegue mvn install su tutti i progetti (usa --tests per abilitare i test)"},
			{"help", "Mostra questa guida"},
		}
		_ = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(tableData).Render()
		pterm.Println()

		// Dettagli comando init
		pterm.DefaultSection.Println("COMANDO: init")
		pterm.FgGray.Println("  Inizializza la configurazione di projman")
		initDetails := []pterm.BulletListItem{
			{Level: 0, Text: "Scansiona la directory specificata alla ricerca di progetti Maven (pom.xml)", Bullet: "•"},
			{Level: 0, Text: "Permette selezione interattiva dei progetti da gestire", Bullet: "•"},
			{Level: 0, Text: "Salva la configurazione in ~/.config/projman/projman_config.json", Bullet: "•"},
		}
		_ = pterm.DefaultBulletList.WithItems(initDetails).Render()
		pterm.Println()

		// Dettagli comando git update
		pterm.DefaultSection.Println("COMANDO: git update")
		pterm.FgGray.Println("  Aggiorna tutti i progetti con gestione intelligente dei branch")
		gitDetails := []pterm.BulletListItem{
			{Level: 0, Text: "Stash automatico delle modifiche non committate", Bullet: "•"},
			{Level: 0, Text: "Selezione interattiva dei progetti da passare a 'develop'", Bullet: "•"},
			{Level: 0, Text: "Pull da origin per branch develop e deploy/*", Bullet: "•"},
			{Level: 0, Text: "Fetch + merge di develop per altri branch", Bullet: "•"},
			{Level: 0, Text: "Ripristino automatico dello stash", Bullet: "•"},
		}
		_ = pterm.DefaultBulletList.WithItems(gitDetails).Render()
		pterm.Println()

		// Dettagli comando mvn install
		pterm.DefaultSection.Println("COMANDO: mvn install")
		pterm.FgGray.Println("  Esegue Maven install su tutti i progetti")
		mvnDetails := []pterm.BulletListItem{
			{Level: 0, Text: "Default: test disabilitati (-DskipTests=true)", Bullet: "•"},
			{Level: 0, Text: "Usa --tests o -t per abilitare l'esecuzione dei test", Bullet: "•"},
			{Level: 0, Text: "Esegue in sequenza su tutti i progetti selezionati", Bullet: "•"},
		}
		_ = pterm.DefaultBulletList.WithItems(mvnDetails).Render()
		pterm.Println()

		// Esempi di utilizzo
		pterm.DefaultSection.Println("ESEMPI")
		examples := []pterm.BulletListItem{
			{Level: 0, Text: "projman init ~/progetti", TextStyle: pterm.NewStyle(pterm.FgLightGreen), Bullet: "→"},
			{Level: 1, Text: "Inizializza projman con i progetti in ~/progetti", Bullet: " "},
			{Level: 0, Text: "projman git update", TextStyle: pterm.NewStyle(pterm.FgLightGreen), Bullet: "→"},
			{Level: 1, Text: "Aggiorna tutti i progetti selezionati", Bullet: " "},
			{Level: 0, Text: "projman mvn install", TextStyle: pterm.NewStyle(pterm.FgLightGreen), Bullet: "→"},
			{Level: 1, Text: "Build di tutti i progetti senza test", Bullet: " "},
			{Level: 0, Text: "projman mvn install --tests", TextStyle: pterm.NewStyle(pterm.FgLightGreen), Bullet: "→"},
			{Level: 1, Text: "Build di tutti i progetti con test abilitati", Bullet: " "},
		}
		_ = pterm.DefaultBulletList.WithItems(examples).Render()
		pterm.Println()

		// Configurazione
		pterm.DefaultSection.Println("CONFIGURAZIONE")
		pterm.DefaultBasicText.Println("La configurazione viene salvata in formato JSON:")
		configPaths := []pterm.BulletListItem{
			{Level: 0, Text: "Windows: %APPDATA%/projman/projman_config.json", Bullet: "📁"},
			{Level: 0, Text: "Linux/macOS: ~/.config/projman/projman_config.json", Bullet: "📁"},
		}
		_ = pterm.DefaultBulletList.WithItems(configPaths).Render()
		pterm.Println()

		// Requisiti
		pterm.DefaultSection.Println("REQUISITI")
		requirements := []pterm.BulletListItem{
			{Level: 0, Text: "Git installato e nel PATH", Bullet: "✓"},
			{Level: 0, Text: "Maven installato e nel PATH (per comandi mvn)", Bullet: "✓"},
			{Level: 0, Text: "Progetti Maven con file pom.xml", Bullet: "✓"},
		}
		_ = pterm.DefaultBulletList.WithItems(requirements).Render()
		pterm.Println()

		// Footer con informazioni aggiuntive
		pterm.DefaultBox.WithTitle("INFO").WithTitleTopCenter().WithBoxStyle(pterm.NewStyle(pterm.FgLightBlue)).Println(
			"Per maggiori informazioni, visita il README.md del progetto\n" +
				"Copyright © 2025 Salvatore Spagnuolo\n" +
				"Licenza: MIT",
		)
	},
}

func init() {
	RootCmd.AddCommand(helpCmd)
}

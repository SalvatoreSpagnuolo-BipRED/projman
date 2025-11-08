/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
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
		// Titolo
		_ = pterm.DefaultHeader.WithFullWidth().Sprint("Projman - Gestore Multiplo di Progetti Git e Maven")
		pterm.DefaultHeader.WithFullWidth().Println("Projman - Gestore Multiplo di Progetti Git e Maven")

		// Uso generale in box
		pterm.DefaultBox.WithTitle("USO").WithPadding(1).WithTitleTopRight().Println("projman [comando] [argomenti] [flags]")

		// Tabella comandi principali
		tableData := pterm.TableData{
			{"COMANDO", "DESCRIZIONE"},
			{"init [directory]", "Scansiona la directory e seleziona i progetti Maven da gestire"},
			{"git update", "Esegue git pull su tutti i progetti selezionati"},
			{"mvn install", "Esegue mvn install su tutti i progetti (usa --tests per abilitarli)"},
		}
		pterm.DefaultTable.WithHasHeader(true).WithData(tableData).Render()

		pterm.Println()
		pterm.DefaultSection.Println("ESEMPI")
		examples := []pterm.BulletListItem{
			{Level: 0, Text: "projman init /path/to/projects", Bullet: "•"},
			{Level: 0, Text: "projman git update", Bullet: "•"},
			{Level: 0, Text: "projman mvn install", Bullet: "•"},
			{Level: 0, Text: "projman mvn install --tests", Bullet: "•"},
		}
		pterm.DefaultBulletList.WithItems(examples).Render()

		pterm.Println()
		pterm.DefaultSection.Println("CONFIGURAZIONE")
		pterm.DefaultBasicText.Println("La configurazione viene salvata in:")
		pterm.FgGray.Println("  • Windows: %APPDATA%/projman/projman_config.json")
		pterm.FgGray.Println("  • Linux/macOS: ~/.config/projman/projman_config.json")
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}

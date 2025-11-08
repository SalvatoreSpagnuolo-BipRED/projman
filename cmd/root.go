package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "projman",
	Short: "Gestione multipla di progetti Git e Maven",
	Long: `Projman è uno strumento da linea di comando che permette di gestire multiple repository Git 
e progetti Maven contemporaneamente. Consente di selezionare un gruppo di progetti 
e eseguire operazioni batch come git pull o mvn install su tutti i progetti selezionati.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		toggle, _ := cmd.Flags().GetBool("toggle")

		fmt.Println("Hello World")
		if toggle {
			fmt.Println("Toggle is true")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

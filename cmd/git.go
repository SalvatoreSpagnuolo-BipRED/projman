/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Esegue git pull su tutti i progetti selezionati",

	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Git Command")
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

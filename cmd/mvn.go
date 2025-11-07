/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// mvnCmd represents the mvn command
var mvnCmd = &cobra.Command{
	Use:   "mvn",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Maven Command")
	},
}

func init() {
	rootCmd.AddCommand(mvnCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mvnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mvnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

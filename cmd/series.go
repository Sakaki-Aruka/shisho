/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Sakaki-Aruka/shisho/cmd/series"
	"github.com/spf13/cobra"
)

// seriesCmd represents the series command
var seriesCmd = &cobra.Command{
	Use:   "series",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("series called")
	},
}



func init() {
	rootCmd.AddCommand(seriesCmd)
	
	seriesCmd.AddCommand(series.SeriesListCmd())
}

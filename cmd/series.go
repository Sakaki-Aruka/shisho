/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Sakaki-Aruka/shisho/cmd/series"
	"github.com/spf13/cobra"
)

// seriesCmd represents the series command
var seriesCmd = &cobra.Command{
	Use:   "series",
	Short: `Series modify commands`,
}

func init() {
	rootCmd.AddCommand(seriesCmd)

	seriesCmd.AddCommand(series.SeriesListCmd())
	seriesCmd.AddCommand(series.SeriesAddCmd())
	seriesCmd.AddCommand(series.SeriesDeleteCmd())
	seriesCmd.AddCommand(series.SeriesModifyCmd())
}

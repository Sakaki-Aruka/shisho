/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `Create BookDB file if it not exists`,
	RunE: cmd,
	// DB がまだ存在しないことを前提とするコマンドのため、
	// rootCmd の PersistentPreRunE (DB 初期化) の伝播を止める。
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
}

func cmd(cmd *cobra.Command, args []string) error {
	if result, err := shishodb.CreateDBFile(dbPath); err != nil {
		return err
	} else {
		if result {
			cmd.Println("BookDB created: " + dbPath)
		} else {
			cmd.Println("BookDB already exists: " + dbPath)
		}
	}

	if err := shishodb.Init(dbPath); err != nil {
		return err
	}

	if err := shishodb.InitTable(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}

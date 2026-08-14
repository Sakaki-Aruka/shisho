/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete registered books",
	Long: `Delete registered books.
Use '--id' or some identifiers (params, k* params) to select targets.`,
	RunE: delete,
}

func delete(cmd *cobra.Command, args []string) error {
	if len(ids) > 0 {
		if err := shishodb.DeleteBooksByIds(ids); err != nil {
			return err
		}
		return nil
	}

	exacts := make(map[string]string)
	parshals := make(map[string]string)
	for key, _ := range searchKeys {
		if s, err := cmd.Flags().GetString(key); err == nil {
			exacts[key] = s
		}

		if k, err := cmd.Flags().GetString("k" + key); err == nil {
			parshals[key] = k
		}
	}

	books, err := shishodb.FindBookByFilter(exacts, parshals)
	if err != nil {
		return err
	}
	
	for _, book := range books {
		ids = append(ids, book.Id)
	}

	if err := shishodb.DeleteBooksByIds(ids); err != nil {
		return err
	}
	return nil
}

var ids []int64

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().Int64SliceVar(&ids, "id", []int64{}, "Delete book ids")
	for k, v := range searchKeys {
		deleteCmd.Flags().String(k, "", v)
		deleteCmd.Flags().String("k" + k, "", v + " keyword")
	}
}

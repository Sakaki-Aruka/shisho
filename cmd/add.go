/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"log"

	"github.com/Sakaki-Aruka/shisho/models"
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a book to the DB.",
	Long: `Add a book to the DB`,
	RunE: add,
}

func add(cmd *cobra.Command, args []string) error {
	argErrors := make([]error, 8)

	author, err := GetFlagOrNil(cmd, "author", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	title, err := GetFlagOrNil(cmd, "title", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	price, err := GetFlagOrNil(cmd, "price", cmd.Flags().GetInt32)
	if err != nil { argErrors = append(argErrors, err) }

	publisher, err := GetFlagOrNil(cmd, "publisher", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	genre, err := GetFlagOrNil(cmd, "genre", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	memo, err := GetFlagOrNil(cmd, "memo", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	isbn, err := GetFlagOrNil(cmd, "isbn", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	janCode, err := GetFlagOrNil(cmd, "jan_code", cmd.Flags().GetString)
	if err != nil { argErrors = append(argErrors, err) }

	actualErrorCount := 0
	for _, e := range argErrors {
		if e != nil && !errors.Is(&pflag.NotExistError{}, e) {
			actualErrorCount++
		}
	}
	if actualErrorCount > 0 {
		return errors.Join(argErrors...)
	}
	
	temp := models.NewBook {
		Author: author,
		Title: title,
		Price: price,
		Publisher: publisher,
		Genre: genre,
		Memo: memo,
		Isbn: isbn,
		JanCode: janCode,
	}

	if useApiDataFiller {
		if err := (&temp).FillDataWithOpenBD(); err != nil {
			return err
		}
	}

	if !temp.CanRegister() {
		return errors.New("Input book data is invalid")
	}

	if !permitDuplicate {
		if err := shishodb.CheckDuplicate(temp); err != nil {
			return err
		}
	}

	if err := shishodb.InsertBook(temp); err != nil {
		log.Default().Printf("Insert error")
		return err
	}

	return nil
}

var permitDuplicate bool
var useApiDataFiller bool

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().String("author", "", "Author name")
	addCmd.Flags().String("title", "", "Book title")
	addCmd.Flags().Int32("price", 0, "Book price")
	addCmd.Flags().String("publisher", "", "Publisher name")
	addCmd.Flags().String("genre", "", "Genre code")
	addCmd.Flags().String("memo", "", "Memo")
	addCmd.Flags().String("isbn", "", "ISBN")
	addCmd.Flags().String("jancode", "", "Jancode (second line)")

	addCmd.Flags().BoolVar(&permitDuplicate, "permit-duplicate", false, "Permit duplicate")
	addCmd.Flags().BoolVar(&useApiDataFiller, "use-filler", true, "Use the book API (OpenBD)")
}

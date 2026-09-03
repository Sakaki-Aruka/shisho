/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/Sakaki-Aruka/shisho/models"
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a book to the DB.",
	Long:  `Add a book to the DB`,
	RunE:  add,
}

func buildNewBookFromFlags(cmd *cobra.Command) (models.NewBook, error) {
	var errs []error

	getOrNil := func(name string) *string {
		v, err := GetFlagOrNil(cmd, name, cmd.Flags().GetString)
		if err != nil {
			errs = append(errs, err)
			return nil
		}
		return v
	}

	author := getOrNil("author")
	title := getOrNil("title")
	publisher := getOrNil("publisher")
	genre := getOrNil("genre")
	memo := getOrNil("memo")
	isbn := getOrNil("isbn")
	janCode := getOrNil("jan_code")

	price, err := GetFlagOrNil(cmd, "price", cmd.Flags().GetInt32)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return models.NewBook{}, errors.Join(errs...)
	}

	return models.NewBook{
		Author:    author,
		Title:     title,
		Price:     price,
		Publisher: publisher,
		Genre:     genre,
		Memo:      memo,
		Isbn:      isbn,
		JanCode:   janCode,
	}, nil
}

func add(cmd *cobra.Command, args []string) error {
	book, err := buildNewBookFromFlags(cmd)
	if err != nil {
		return err
	}

	if useApiDataFiller {
		if err := (&book).FillDataWithOpenBD(); err != nil {
			return err
		}
	}

	if !book.CanRegister() {
		return errors.New("Input book data is invalid")
	}

	if !permitDuplicate {
		if err := shishodb.CheckDuplicate(book); err != nil {
			return err
		}
	}

	if err := shishodb.InsertBook(book); err != nil {
		log.Default().Printf("Insert error")
		return err
	}

	fmt.Printf("Success to add a book. (Title: %s, Isbn: %s)\n", *book.Title, *book.Isbn)

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
	addCmd.Flags().String("jan_code", "", "Jancode (second line)")

	addCmd.Flags().BoolVar(&permitDuplicate, "permit-duplicate", false, "Permit duplicate")
	addCmd.Flags().BoolVar(&useApiDataFiller, "use-filler", true, "Use the book API (OpenBD)")
}

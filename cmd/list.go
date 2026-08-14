/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Sakaki-Aruka/shisho/models"
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  ``,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("%d books found.\n", len(books))
	for _, book := range books {
		line := ""

		if lineformat != "" {
			line, err = formatted(book)
			if err != nil {
				return err
			}
		}

		if lineformat == "" || line == "" {
			continue
		}

		fmt.Println(line)
	}
	return nil
}

func formatted(book models.Book) (string, error) {
	t, err := template.New("").Parse(lineformat)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err := t.Execute(buf, book.ToMap()); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var searchKeys = map[string]string{
	// key, short description
	"author":    "Author",
	"title":     "Title",
	"publisher": "Publisher",
	"genre":     "Genre (C code)",
	"memo":      "Memo",
	"isbn":      "ISBN code",
	"jancode":   "Book JAN code (second line)",
}

var lineformat string

func init() {
	// shishodb.FindBookByFilter で空文字列のフィルタは排除されるので、デフォルトを空文字列にしておけばそのまま検索用フィルタに利用できる
	rootCmd.AddCommand(listCmd)

	for key, usage := range searchKeys {
		listCmd.Flags().String(key, "", usage)
		listCmd.Flags().String("k"+key, "", usage+" keyword")
	}

	listCmd.Flags().StringVar(&lineformat, "format", "Id: {{ printf \"%04d\" .Id }}, Title: {{.Title}}, Author: {{.Author}}, ISBN: {{.ISBN}}", "Output line format")
}

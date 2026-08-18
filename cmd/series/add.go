package series

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sakaki-Aruka/shisho/code"
	"github.com/Sakaki-Aruka/shisho/models"
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

type seriesAddOption struct {
	Title string
	Stauts models.Status
	Isbns []string
}

func (opt *seriesAddOption) isValid() error {
	if opt.Title == "" {
		return errors.New("Title must not be empty.")
	}

	invalid := make([]string, 0)
	for _, isbn := range opt.Isbns {
		if !code.IsValidIsbn(isbn) {
			invalid = append(invalid, isbn)
		}
	}

	if len(invalid) > 0 {
		return fmt.Errorf("Invalid isbn found (%d): %s", len(invalid), strings.Join(invalid, ", "))
	}

	return nil
}

func SeriesAddCmd() *cobra.Command {
	opt := &seriesAddOption{}

	cmd := &cobra.Command{
		Use: "add",
		Short: "Add series",
		RunE: func(cmd *cobra.Command, args []string) error { 
			return seriesAdd(opt)
		},
	}

	cmd.Flags().StringVar(&opt.Title, "title", "", "Series title")
	cmd.Flags().Var(&opt.Stauts, "status", "Series status (ongoing / completed)")
	cmd.Flags().StringSliceVar(&opt.Isbns, "isbn", []string{}, "Series elements isbn list")

	return cmd
}

func seriesAdd(opt *seriesAddOption) error {
	if err := opt.isValid(); err != nil {
		return err
	}

	new := models.NewSeries {
		Title: opt.Title,
		Status: opt.Stauts,
		Isbns: opt.Isbns,
	}

	if err := shishodb.InsertSeries(&new); err != nil {
		return err
	}

	fmt.Println("Success to add series")

	return nil
}
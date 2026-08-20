package series

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

func SeriesModifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify",
		Short: "Modify series data",
	}

	cmd.AddCommand(seriesModifyTitle())
	cmd.AddCommand(seriesModifyDeleteIsbn())
	cmd.AddCommand(seriesModifyAddIsbn())

	return cmd
}

func seriesModifyTitle() *cobra.Command {
	return &cobra.Command{
		Use:   "title [id] [title]",
		Args:  cobra.ExactArgs(2),
		Short: "Modify series title",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.New("'id' parse error")
			}
			new := args[1]
			if new == "" {
				return errors.New("'title' must not be empty")
			}

			if err := shishodb.UpdateSeriesTitle(id, new); err != nil {
				return err
			}
			fmt.Println("Success to update series title")
			return nil
		},
	}
}

func seriesModifyDeleteIsbn() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-isbn [id] [isbn...]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Delete isbn from the specified series",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.New("Series id parse error. It must be integer")
			}
			isbns := args[1:]
			isbns = slices.DeleteFunc(isbns, func(s string) bool { return s == "" })
			if len(isbns) == 0 {
				return errors.New("isbn must not be empty")
			}

			affectedRows, err := shishodb.DeleteSeriesIsbn(id, &isbns)
			if err != nil {
				return err
			}

			affected := ""
			if affectedRows > -1 {
				affected = strconv.Itoa(int(affectedRows))
			} else {
				affected = "?"
			}
			
			fmt.Printf("Success to delete  isbn(s) from the specified series (affected: %s)\n", affected)
			
			return nil
		},
	}
}

func seriesModifyAddIsbn() *cobra.Command {
	return &cobra.Command{
		Use: "add-isbn [id] [isbn...]",
		Args: cobra.MinimumNArgs(2),
		Short: "Add isbn to the specified series",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.New("Series id parse error. It must be integer")
			}
			isbns := args[1:]
			isbns = slices.DeleteFunc(isbns, func(s string) bool { return s == "" })
			if len(isbns) == 0 {
				return errors.New("isbn must not be empty")
			}

			affected, err := shishodb.AddSeriesIsbn(id, &isbns)
			if err != nil {
				return err
			}

			if affected == 0 {
				fmt.Println("No errors occurred, but isbn(s) not inserted")
			} else {
				fmt.Printf("Isbns inserted (%d)\n", affected)
			}

			return nil
		},
	}
}

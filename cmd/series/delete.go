package series

import (
	"fmt"
	"strconv"
	"strings"

	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

type seriesDeleteOption struct {
	Ids []int64
}

func SeriesDeleteCmd() *cobra.Command {
	opt := &seriesDeleteOption{}

	cmd := &cobra.Command{
		Use: "delete",
		Short: "Delete series",
		RunE: func(cmd *cobra.Command, args []string) error {
			return seriesDelete(opt)
		},
	}

	cmd.Flags().Int64SliceVar(&opt.Ids, "id", []int64{}, "Delete series ids")

	return cmd
}

func seriesDelete(opt *seriesDeleteOption) error {
	if len(opt.Ids) == 0 {
		fmt.Println("Id not specified")
		return nil
	}

	success := make([]string, 0)
	for _, id := range opt.Ids {
		affected, err := shishodb.DeleteSeriesById(id);
		if err != nil {
			return err
		}

		if affected > 0 {
			success = append(success, strconv.FormatInt(id, 10))
		}
	}

	if len(success) > 0 {
		fmt.Printf("Success to delete series (id: %s)\n", strings.Join(success, ", "))
	} else {
		fmt.Println("Series not deleted")
	}

	return nil
}
package series

import (
	"fmt"
	"strconv"
	"strings"

	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

func SeriesDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete [id...]",
		Args: cobra.MinimumNArgs(1),
		Short: "Delete series",
		RunE: func(cmd *cobra.Command, args []string) error {
			var ids []int64
			invalids := make([]string, 0)
			for _, id := range args {
				if n, err := strconv.ParseInt(id, 10, 64); err != nil {
					invalids = append(invalids, id)
				} else {
					ids = append(ids, n)
				}
			}

			if len(invalids) != 0 {
				return fmt.Errorf("Id must be integer (invalid: %s)\n", strings.Join(invalids, ", "))
			}
			return seriesDelete(&ids)
		},
	}
	return cmd
}

func seriesDelete(ids *[]int64) error {
	if len(*ids) == 0 {
		fmt.Println("Id not specified")
		return nil
	}

	success := make([]string, 0)
	for _, id := range *ids {
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
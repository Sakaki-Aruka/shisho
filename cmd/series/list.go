package series

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"text/template"

	"github.com/Sakaki-Aruka/shisho/models"
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

type seriesListOption struct {
	Id      int64
	Title   string
	KTitle  string
	Json    bool
	Jsonl   bool
	Oneline bool
	Not     bool
}

func (opt seriesListOption) isValid() error {
	if opt.Json && opt.Jsonl {
		return errors.New("Invalid parameter combination (--json and --jsonl)")
	}

	if (opt.Json || opt.Jsonl) && opt.Oneline {
		return errors.New("Invalid parameter combination ({--json or --jsonl} and --1)")
	}

	return nil
}

func SeriesListCmd() *cobra.Command {
	option := &seriesListOption{}

	cmd := &cobra.Command{
		Use: "list",
		Short: "Displays specified serieses",
		RunE: func(cmd *cobra.Command, args []string) error {
			return seriesList(option)
		},
	}

	cmd.Flags().Int64Var(&option.Id, "id", -1, "Series id")
	cmd.Flags().StringVar(&option.Title, "title", "", "Title")
	cmd.Flags().StringVar(&option.KTitle, "ktitle", "", "Title keyword")
	cmd.Flags().BoolVar(&option.Json, "json", false, "Use json format")
	cmd.Flags().BoolVar(&option.Jsonl, "jsonl", false, "Use jsonl format")
	cmd.Flags().BoolVar(&option.Oneline, "1", false, "1 ISBN / 1 line")
	cmd.Flags().BoolVar(&option.Not, "not", false, "")

	return cmd
}

func seriesList(opt *seriesListOption) error {
	if err := opt.isValid(); err != nil {
		return err
	}

	exacts := make(map[string]string)
	parshals := make(map[string]string)

	if opt.Id != -1 {
		exacts["id"] = strconv.Itoa(int(opt.Id))
	}

	if opt.Title != "" {
		exacts["title"] = opt.Title
	}

	if opt.KTitle != "" {
		parshals["title"] = opt.KTitle
	}

	series, err := shishodb.FindSeriesByFilter(exacts, parshals)
	if err != nil {
		fmt.Println("Failed to find specified series")
		return err
	} else if len(series) == 0 {
		fmt.Println("No series found")
		return nil
	}

	if opt.Not {
		for i := range series {
			if len(series[i].Isbns) == 0 {
				continue
			}
			if err := removePosessions(&series[i]); err != nil {
				return err
			}
		}

		series = slices.DeleteFunc(series, func(s models.Series) bool { return len(s.Isbns) == 0 })
	}

	if len(series) == 0 {
			return nil
		}

	if opt.Json {
		fmt.Println(formatJson(&series))
	} else if opt.Jsonl {
		for _, line := range formatJsonl(&series) {
			fmt.Println(line)
		}
	} else if opt.Oneline {
		if lines, err := oneline(&series); err == nil {
			for _, line := range lines {
				fmt.Println(line)
			}
		}
	} else {
		for _, line := range shortLine(&series) {
			fmt.Println(line)
		}
	}

	return nil
}

func shortLine(series *[]models.Series) []string {
	var result []string
	t, _ := template.New("").Parse("Id: {{printf \"%04d\" .Id}}, Filter matched: {{.Amount}}, Title: {{.Title}}")
	buf := bytes.NewBufferString("")
	for _, s := range *series {
		values := struct { 
			Id int64
			Title string
			Amount int
		} { 
			Id: s.Id,
			Title: s.Title,
			Amount: len(s.Isbns),
		}

		if err := t.Execute(buf, &values); err != nil {
			continue
		}

		result = append(result, buf.String())
		buf.Reset()
	}

	return result
}

func formatJson(series *[]models.Series) string {
	result := ""
	j, err := json.Marshal(series)
	if err == nil {
		result = string(j)
	}
	return result
}

func formatJsonl(series *[]models.Series) []string {
	var result []string
	for _, s := range *series {
		j, err := json.Marshal(s)
		if err != nil {
			continue
		}

		result = append(result, string(j))
	}
	return result
}

func oneline(series *[]models.Series) ([]string, error) {
	var result []string
	for _, s := range *series {
		result = append(result, fmt.Sprintf("[%d] Title: %s, (%v)", s.Id, s.Title, s.Status))
		result = append(result, s.Isbns...)
	}

	return result, nil
}

func removePosessions(series *models.Series) error {
	trimed, err := shishodb.GetNotOwnIsbn(series.Isbns)
	if err != nil {
		return err
	}
	series.Isbns = trimed
	return nil
}

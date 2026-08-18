package shishodb

import (
	"errors"
	"fmt"
	"strings"

	models "github.com/Sakaki-Aruka/shisho/models"
)

func InsertSeries(n *models.NewSeries) error {
	result, err := db.Exec("INSERT INTO series (title, status) VALUES (?, ?)", n.Title, n.Status.String())
	if err != nil {
		return err
	}

	if len(n.Isbns) > 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return errors.New("Series Isbn insert error (generated id not found)")
		}
		for i := range n.Isbns {
			isbn := n.Isbns[i]
			if err := InsertIsbnToSeries(id, isbn); err != nil {
				return errors.New("Series Isbn insert error (insert isbn error)")
			}
		}
	}
	return nil
}

func InsertIsbnToSeries(id int64, isbn string) error {
	if _, err := db.Exec("INSERT INTO series_isbn VALUES (?, ?)", id, isbn); err != nil {
		return err
	}

	return nil
}

func FindSeriesById(id int64) (*models.Series, error) {
	var series models.Series
	if err := db.Select(&series, "SELECT * FROM series WHERE id = ?", id); err != nil {
		return nil, err
	}

	if series.Id == 0 {
		return nil, fmt.Errorf("[Error] Series not found. (Id: %d)", id)
	}

	isbns, err := FindIsbnsBySeriesId(series.Id)
	if err != nil {
		return nil, err
	}
	series.Isbns = isbns

	return &series, nil
}

func FindIsbnsBySeriesId(id int64) ([]string, error) {
	var isbns []string
	if err := db.Select(&isbns, "SELECT isbn FROM series_isbn WHERE series_id = ? ORDER BY isbn", id); err != nil {
		return nil, err
	}

	return isbns, nil
}

func FindSeriesByFilter(exactFilters map[string]string, parshalFilters map[string]string) ([]models.Series, error) {
	var filters []string
	if len(exactFilters) > 0 {
		for k, v := range exactFilters {
			if v != "" {
				filters = append(filters, buildExactFilter(k, v))
			}
		}
	}

	if len(parshalFilters) > 0 {
		for k, v := range parshalFilters {
			if v != "" {
				filters = append(filters, buildParshalFilter(k, v))
			}
		}
	}

	where := strings.Join(filters, " AND ")
	sql := "SELECT * FROM series"
	if len(filters) > 0 {
		sql = sql + " WHERE " + where
	}

	var series []models.Series
	if err := db.Select(&series, sql); err != nil {
		return nil, err
	}

	for i := 0; i < len(series); i++ {
		s := &series[i]
		isbns, err := FindIsbnsBySeriesId(s.Id)
		if err == nil {
			s.Isbns = isbns
		}
	}

	return series, nil
}

func DeleteSeriesById(id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("Delete series error (id: %d)", id)
	}

	affected, _ := result.RowsAffected()

	if _, err := db.Exec("DELETE FROM series_isbn WHERE series_id = ?", id); err != nil {
		return affected, fmt.Errorf("Delete series isbn error (series id: %d)", id)
	}

	return affected, nil
}

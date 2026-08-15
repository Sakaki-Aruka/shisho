package shishodb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	models "github.com/Sakaki-Aruka/shisho/models"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

func Init(dbPath string) error {
	// empty path string means using the default setting
	var abs string
	if dbPath == "" {
		if p, err := GetDefaultPath(); err != nil {
			return err
		} else {
			abs, err = filepath.Abs(p)
			if err != nil {
				return err
			}
		}
	} else {
		a, err := filepath.Abs(dbPath)
		if err != nil {
			return err
		}
		abs = a
	}

	if d, err := sqlx.Connect("sqlite", abs); err != nil {
		return err
	} else {
		db = d
		return nil
	}
}

func CreateDBFile(path string) (bool, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(abs); err == nil {
		return false, nil
	}

	if err := os.MkdirAll(filepath.Dir(abs), os.ModePerm); err != nil {
		return false, err
	}

	if _, err := os.Create(abs); err != nil {
		return false, err
	}

	return true, nil
}

func InitTable() error {
	books := `
CREATE TABLE IF NOT EXISTS books (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	author TEXT,
	title TEXT NOT NULL,
	price INTEGER,
	publisher TEXT,
	genre TEXT,
	memo TEXT,
	isbn TEXT,
	jan_code TEXT
);`

	series := `
CREATE TABLE IF NOT EXISTS series (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	status TEXT NOT NULL
);`

	seriesIsbn := `
CREATE TABLE IF NOT EXISTS series_isbn (
	series_id INTEGER NOT NULL REFERENCES series(id),
	isbn TEXT NOT NULL,
	PRIMARY KEY (series_id, isbn)
);
`

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	tx.Exec(books)
	tx.Exec(series)
	tx.Exec(seriesIsbn)

	return nil
}

func GetDefaultPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, ".config", "shisho", "book.db")
	return path, nil
}

func FindBookById(id int64) (models.Book, error) {
	book := models.Book{}
	err := db.Get(&book, "SELECT * FROM books WHERE id = ?", id)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func FindAllBooks() ([]models.Book, error) {
	books := []models.Book{}
	err := db.Select(&books, "SELECT * FROM books")
	if err != nil {
		return nil, err
	}

	return books, nil
}

func FindAllIsbn() ([]string, error) {
	var codes []string
	err := db.Select(&codes, "SELECT isbn FROM books WHERE isbn IS NOT NULL")
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func DeleteBookById(id int64) error {
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBooksByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query, params, err := sqlx.In("DELETE FROM books WHERE id IN (?)", ids)
	if err != nil {
		fmt.Println("Batch book delete error")
		return err
	}

	result, err := db.Exec(query, params...)
	if err != nil {
		fmt.Println("Batch book delete error")
		return err
	}

	if affected, err := result.RowsAffected(); err == nil {
		fmt.Printf("%d books deleted\n", affected)
	}

	return nil
}

func InsertBook(book models.NewBook) error {
	sql := `
INSERT INTO books (author, title, price, publisher, genre, memo, isbn, jan_code)
VALUES (:author, :title, :price, :publisher, :genre, :memo, :isbn, :jan_code)
`
	if _, err := db.NamedExec(sql, book.ToBook()); err != nil {
		return err
	}

	return nil
}

func FindBookByFilter(exactFilters map[string]string, parshalFilters map[string]string) ([]models.Book, error) {
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
	sql := "SELECT * FROM books"
	if len(filters) > 0 {
		sql = sql + " WHERE " + where
	}

	//debug
	fmt.Printf("sql: %s, where: %s\n", sql, where)

	var books []models.Book
	if err := db.Select(&books, sql); err != nil {
		return nil, err
	}

	return books, nil
}

func buildExactFilter(columnName string, value string) string {
	return columnName + " LIKE '" + value + "'"
}

func buildParshalFilter(columnName string, value string) string {
	return columnName + " LIKE '%" + value + "%'"
}

func CheckDuplicate(n models.NewBook) error {
	// 重複禁止。識別子 (title と ISBN コードが存在しない場合は検索できないのでエラー扱い)
	exactFilters := make(map[string]string)
	if n.Title != nil {
		exactFilters["title"] = *n.Title
	}
	if n.Isbn != nil {
		exactFilters["isbn"] = *n.Isbn
	}

	if len(exactFilters) == 0 {
		return errors.New("Book identifiers (title and ISBN code) are not found. Specify those.")
	}

	if stores, err := FindBookByFilter(exactFilters, make(map[string]string)); err != nil {
		return err
	} else if len(stores) > 0 {
		return errors.New("Book identifier duplicated. ")
	} else {
		return nil
	}
}

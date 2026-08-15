package shishodb

import (
	"os"
	"path/filepath"

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
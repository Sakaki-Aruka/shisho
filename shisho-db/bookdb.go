package shishodb

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	models "github.com/Sakaki-Aruka/shisho/models"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

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

	var books []models.Book
	if err := db.Select(&books, sql); err != nil {
		return nil, err
	}

	return books, nil
}

func GetNotOwnIsbn(candidate []string) ([]string, error) {
	if len(candidate) == 0 {
		return nil, errors.New("Own books isbn must not empty")
	}

	var owned []string
	q, params, err := sqlx.In("SELECT isbn FROM books WHERE isbn IN (?)", candidate)
	if err != nil {
		return nil, err
	}

	q = db.Rebind(q)
	if err := db.Select(&owned, q, params...); err != nil {
		return nil, err
	}

	ownedSet := make(map[string]struct{}, len(owned))
	for _, isbn := range owned {
		ownedSet[isbn] = struct{}{}
	}

	notOwned := slices.DeleteFunc(slices.Clone(candidate), func(s string) bool {
		_, ok := ownedSet[s]
		return ok
	})

	return notOwned, nil
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

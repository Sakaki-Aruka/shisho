package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Sakaki-Aruka/shisho/code"
	"github.com/tidwall/gjson"
)

type Book struct {
	Id        int64          `db:"id"`
	Author    sql.NullString `db:"author"`
	Title     sql.NullString `db:"title"`
	Price     sql.NullInt32  `db:"price"`
	Publisher sql.NullString `db:"publisher"`
	Genre     sql.NullString `db:"genre"`
	Memo      sql.NullString `db:"memo"`
	Isbn      sql.NullString `db:"isbn"`
	JanCode   sql.NullString `db:"jan_code"`
}

func (b Book) ToMap() map[string]any {
	unwrapStr := func(s sql.NullString) string {
		if s.Valid {
			return s.String
		} else {
			return "<nil>"
		}
	}

	unwrapInt32 := func(i sql.NullInt32) int32 {
		if i.Valid {
			return i.Int32
		} else {
			return -1
		}
	}

	elements := make(map[string]any)
	elements["Id"] = b.Id
	elements["Author"] = unwrapStr(b.Author)
	elements["Title"] = unwrapStr(b.Title)
	elements["Price"] = unwrapInt32(b.Price)
	elements["Publisher"] = unwrapStr(b.Publisher)
	elements["Genre"] = unwrapStr(b.Genre)
	elements["Memo"] = unwrapStr(b.Memo)
	elements["ISBN"] = unwrapStr(b.Isbn)
	elements["JanCode"] = unwrapStr(b.JanCode)
	return elements
}

type NewBook struct {
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Price     *int32  `json:"price"`
	Publisher *string `json:"publisher"`
	Genre     *string `json:"genre"`
	Memo      *string `json:"memo"`
	Isbn      *string `json:"isbn"`
	JanCode   *string `json:"jan_code"`
}

func getOrDefault[T any](v *T) T {
	var result T
	if v == nil {
		return result
	}
	return *v
}

func (n NewBook) ToBook() Book {
	return Book{
		Id:        -1,
		Author:    sql.NullString{String: getOrDefault(n.Author), Valid: n.Author != nil},
		Title:     sql.NullString{String: getOrDefault(n.Title), Valid: n.Title != nil},
		Price:     sql.NullInt32{Int32: getOrDefault(n.Price), Valid: n.Price != nil},
		Publisher: sql.NullString{String: getOrDefault(n.Publisher), Valid: n.Publisher != nil},
		Genre:     sql.NullString{String: getOrDefault(n.Genre), Valid: n.Genre != nil},
		Memo:      sql.NullString{String: getOrDefault(n.Memo), Valid: n.Memo != nil},
		Isbn:      sql.NullString{String: getOrDefault(n.Isbn), Valid: n.Isbn != nil},
		JanCode:   sql.NullString{String: getOrDefault(n.JanCode), Valid: n.JanCode != nil},
	}
}

func (n NewBook) CanRegister() bool {
	return !((n.Title == nil || *n.Title == "") && (n.Isbn == nil || *n.Isbn == "" || !code.IsValidIsbn(*n.Isbn)))
}

func (n *NewBook) FillDataWithOpenBD() error {
	if n.Isbn == nil {
		return errors.New("OpenDB needs ISBN code for searching book data")
	}

	url := fmt.Sprintf("https://api.openbd.jp/v1/get?isbn=%s", *n.Isbn)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid book data API response (%d)", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return errors.New("Book data API response is empty")
	}

	setData(body, n)
	return nil
}

func setData(data []byte, book *NewBook) {
	if author := gjson.GetBytes(data, "0.summary.author").String(); author != "" {
		book.Author = &author
	}

	if title := gjson.GetBytes(data, "0.summary.title").String(); title != "" {
		book.Title = &title
	}

	if price := gjson.GetBytes(data, "0.onix.ProductSupply.SupplyDetail.Price.0.PriceAmount").Int(); price >= 0 {
		p := int32(price)
		book.Price = &p
	}

	if publisher := gjson.GetBytes(data, "0.summary.publisher").String(); publisher != "" {
		book.Publisher = &publisher
	}
}

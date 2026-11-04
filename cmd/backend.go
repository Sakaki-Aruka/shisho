package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sakaki-Aruka/shisho/models"
	shishodb "github.com/Sakaki-Aruka/shisho/shisho-db"
	"github.com/spf13/cobra"
)

// やりたいこと
// REST API からの追加・削除・確認

var handleCmd = &cobra.Command{
	Use:   "handle",
	Short: "Start an api server",
	RunE:  handle,
}

func handle(cmd *cobra.Command, args []string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", indexEndpoint)
	mux.HandleFunc("QUERY /books", bookGetEndpoint)
	mux.HandleFunc("POST /books", bookAddEndpoint)
	mux.HandleFunc("DELETE /books", bookDeleteEndpoint)

	log.Printf("listening on localhost:%d\n", port)
	if token != "" {
		log.Println("token requires")
	}
	err := http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), mux)
	if err != nil {
		return err
	}
	return nil
}

//go:embed web/index.html
var indexHTML []byte

// indexEndpoint serves the web UI (browse / add / delete books).
func indexEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(indexHTML)
}

// decodeJSONBody はリクエストボディを v にデコードする。
// ボディが空なら何もしない (全フィールドがゼロ値のまま)。
// 構文エラーや型の不一致があれば 400 を返して false を返す。
func decodeJSONBody(w http.ResponseWriter, r *http.Request, v any, tag string) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("err (%s): failed to read a request body. (%s)\n", tag, err.Error())
		return false
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return true
	}
	if err := json.Unmarshal(body, v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		log.Printf("err (%s): invalid request body. (%s)\n", tag, err.Error())
		return false
	}
	return true
}

func bookGetEndpoint(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Parshals map[string]string `json:"parshals"`
		Exacsts  map[string]string `json:"exacts"`
		Token    string            `json:"token"`
	}

	if !decodeJSONBody(w, r, &req, "get") {
		return
	}
	if token != "" {
		if req.Token != token {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("err (get): authentication failed")
			return
		}
	}

	books, err := shishodb.FindBookByFilter(req.Exacsts, req.Parshals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("err (get): failed to find all books. (%s)\n", err.Error())
		return
	}
	j, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("err (get): failed to data mashalling. (%s)\n", err.Error())
		return
	}

	log.Printf("info (get): request got %d books.\n", len(books))
	w.Write(j)
}

func bookAddEndpoint(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token          string         `json:"token"`
		CheckDuplicate bool           `json:"check-duplicate"`
		FillBookData   bool           `json:"fill"`
		Params         models.NewBook `json:"params"`
	}

	if !decodeJSONBody(w, r, &req, "add") {
		return
	}
	if token != "" {
		if req.Token != token {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("err (add): authentication failed")
			return
		}
	}

	if req.CheckDuplicate {
		if err := shishodb.CheckDuplicate(req.Params); err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("book duplicated."))
			log.Println("err (add): duplicated")
			return
		}
	}

	if !req.Params.CanRegister() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid 'params'"))
		log.Printf("err (add): invalid params.")
		return
	}

	if req.FillBookData {
		if err := (&req.Params).FillDataWithOpenBD(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to fill data"))
			log.Printf("err (add): fill error. (%s)\n", err.Error())
			return
		}
	}

	if err := shishodb.InsertBook(req.Params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("err (add): insert error. (%s)\n", err.Error())
		return
	}

	book := req.Params.ToBook()
	j, _ := json.Marshal(book)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("info (add): book added (%s)\n", string(j))
}

func bookDeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string  `json:"token"`
		Ids   []int64 `json:"ids"`
	}

	if !decodeJSONBody(w, r, &req, "del") {
		return
	}
	if token != "" {
		if req.Token != token {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("err (del): authentication failed")
			return
		}
	}

	if len(req.Ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("need some ids"))
		log.Println("err (del): zero ids")
		return
	}

	if err := shishodb.DeleteBooksByIds(req.Ids); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("err (del): failed to delete. (%s)\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	strlist := []string{}
	for _, id := range req.Ids {
		strlist = append(strlist, strconv.FormatInt(id, 10))
	}
	log.Printf("info (del): delete books. (ids: %s)\n", strings.Join(strlist, ", "))
}

var port uint16
var token string

func init() {
	rootCmd.AddCommand(handleCmd)
	handleCmd.Flags().Uint16Var(&port, "port", 20215, "Server port")
	handleCmd.Flags().StringVar(&token, "token", "", "Authenticate token")
}

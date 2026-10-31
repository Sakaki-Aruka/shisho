package cmd

import (
	"encoding/json"
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

	mux.HandleFunc("GET /books", bookGetEndpoint)
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

func bookGetEndpoint(w http.ResponseWriter, r *http.Request) {
	if token != "" {
		if t := r.Header.Get("token"); t != token {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("err (get): authentication failed")
			return
		}
	}
	parshals := map[string]string{}
	exacts := map[string]string{}
	if p := r.Header.Get("parshals"); p != "" {
		json.Unmarshal([]byte(p), &parshals)
	}
	if ex := r.Header.Get("exacts"); ex != "" {
		json.Unmarshal([]byte(ex), &exacts)
	}

	books, err := shishodb.FindBookByFilter(exacts, parshals)
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
	if token != "" {
		if t := r.Header.Get("token"); t != token {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("err (get): authentication failed")
			return
		}
	}
	
	params := r.Header.Get("params")
	if params == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("need 'params'"))
		log.Println("err (add): no params")
		return
	}

	new := models.NewBook{}
	if err := json.Unmarshal([]byte(params), &new); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid 'params'"))
		log.Printf("err (add): invalid params. (%s)\n", err.Error())
		return
	}

	fill := r.Header.Get("fill")
	if strings.ToLower(fill) == "true" {
		if err := (&new).FillDataWithOpenBD(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to fill data"))
			log.Printf("err (add): fill error. (%s)\n", err.Error())
			return
		}
	}

	if err := shishodb.InsertBook(new); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("err (add): insert error. (%s)\n", err.Error())
		return
	}

	book := new.ToBook()
	j, _ := json.Marshal(book)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("info (add): book added (%s)\n", string(j))
}

func bookDeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	if token != "" {
		if t := r.Header.Get("token"); t != token {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("err (get): authentication failed")
			return
		}
	}

	ids := []int64{}
	i := r.Header.Get("ids")
	if i == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("need the 'ids' header"))
		log.Println("err (del): no 'ids'")
		return
	}

	if err := json.Unmarshal([]byte(i), &ids); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid 'ids'"))
		log.Printf("err (del): invalid ids. (%s)\n", err.Error())
		return
	}

	if err := shishodb.DeleteBooksByIds(ids); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("err (del): failed to delete. (%s)\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	strlist := []string{}
	for _, id := range ids {
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

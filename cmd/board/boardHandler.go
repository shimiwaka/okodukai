package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"github.com/go-chi/chi"
)

func boardHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))


	if board.Owner == "" {
		fmt.Fprintln(w, "{\"error\": \"invalid token\"}")
	} else {
		columns := []schema.Column{}
		db.Where("board = ?", board.Token).Find(&columns)

		resp := schema.Response{
			Owner: board.Owner,
			Token: board.Token,
			Columns: columns,
		}
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "{\"error\": true, \"message\": \"failed to encode json\"}")
			return
		}
		fmt.Fprint(w, buf.String())
	}
	db.Close()
}

func addColumnHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseForm()
	if e != nil {
		panic("error: parse error occured.")
	}

	name := r.Form.Get("name")

	db := connector.ConnectDB()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	column := schema.Column{Board: board.Token, Name: name}

	db.Create(&column)
	db.Close()
}

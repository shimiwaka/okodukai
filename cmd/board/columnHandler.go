package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"github.com/go-chi/chi"
)

func addColumnHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseForm()
	if e != nil {
		fmt.Fprintf(w, "{\"success\":false, \"message\":\"parse error occured\"}")
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))
	
	if name == "" || price < 0 {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid parameter\"}")
		return
	}

	db := connector.ConnectDB()
	defer db.Close()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	if board.Owner == "" {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid token\"}")
		return
	}

	column := schema.Column{Board: board.ID, Name: name, Price: price}

	db.Create(&column)
}

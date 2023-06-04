package main

import (
	"fmt"
	"time"
	"net/http"
	"strconv"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"github.com/go-chi/chi"
)

func checkHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	if board.Owner == "" {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid token\"}")
		return
	}

	columns := []schema.Column{}
	db.Where("board = ?", board.ID).Order("id asc").Find(&columns)

	columnIdx, _ := strconv.Atoi(chi.URLParam(r, "column"))

	if columnIdx > len(columns) || columnIdx < 0 {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid column number\"}")
		return
	}
	column := columns[columnIdx]

	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", chi.URLParam(r, "date")))
	t = t.Add(time.Hour * -9)
	check := schema.Check{Date: t, Column: column.ID, Board: board.ID}

	db.Create(&check)
}


func uncheckHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	if board.Owner == "" {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid token\"}")
		return
	}

	columns := []schema.Column{}
	db.Where("board = ?", board.ID).Find(&columns)

	columnIdx, _ := strconv.Atoi(chi.URLParam(r, "column"))
	if columnIdx > len(columns) || columnIdx < 0 {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid column number\"}")
		return
	}
	column := columns[columnIdx]

	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", chi.URLParam(r, "date")))
	t = t.Add(time.Hour * -9)

	db.Where("board = ? AND `column` = ? AND date = ?", board.ID, column.ID, t).Delete(&schema.Check{})
}

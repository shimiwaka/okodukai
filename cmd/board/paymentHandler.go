package main

import (
	"fmt"
	"time"
	"net/http"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"github.com/go-chi/chi"
)

func newPaymentHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	if board.Owner == "" {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid token\"}")
		return
	}

	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", chi.URLParam(r, "date")))
	t = t.Add(time.Hour * -9)
	payment := schema.Payment{Date: t, Board: board.ID}

	db.Create(&payment)
	fmt.Fprintln(w, "{\"success\": true}")
}

func cancelPaymentHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()
	defer db.Close()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	if board.Owner == "" {
		fmt.Fprintln(w, "{\"success\": false, \"message\": \"invalid token\"}")
		return
	}

	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", chi.URLParam(r, "date")))
	t = t.Add(time.Hour * -9)

	db.Where("board = ? AND date = ?", board.ID, t).Delete(&schema.Payment{})
	fmt.Fprintln(w, "{\"success\": true}")
}
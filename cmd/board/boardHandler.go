package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"net/http"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"github.com/go-chi/chi"
)

func TrancateByDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func showBoardHandler(w http.ResponseWriter, r *http.Request) {
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

	createdAt := TrancateByDate(board.CreatedAt.AddDate(0, 0, -5))
	now := TrancateByDate(time.Now()).AddDate(0, 0, 1)

	days := []schema.Day{}

	checks := []schema.Check{}
	db.Where("board = ?", board.ID).Find(&checks)
	checkedMap := make(map[string]bool)

	for _, v := range(checks) {
		checkedMap[fmt.Sprintf("%s_%d", v.Date, v.Column)] = true
	}

	for i := now; i.After(createdAt); i = i.AddDate(0, 0, -1) {
		checkedList := []bool{}
		for _, column := range(columns){
			if checkedMap[fmt.Sprintf("%s_%d", i, column.ID)] {
				checkedList = append(checkedList, true)
			} else {
				checkedList = append(checkedList, false)
			}
		}
		days = append(days, schema.Day{Date: i, Checked: checkedList})
	}

	resp := schema.Response{
		Owner: board.Owner,
		Token: board.Token,
		Columns: columns,
		Days: days,
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "{\"success\": false, \"message\": \"failed to encode json\"}")
		return
	}
	fmt.Fprint(w, buf.String())
}
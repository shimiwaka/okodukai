package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"net/http"
	"strconv"

	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	"github.com/go-chi/chi"
)

func TrancateByDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func boardHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	if board.Owner == "" {
		fmt.Fprintln(w, "{\"error\": \"invalid token\"}")
	} else {
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

		for i := createdAt; i.Before(now); i = i.AddDate(0, 0, 1) {
			checkedList := []bool{}
			for j := 0; j < len(columns); j++ {
				if checkedMap[fmt.Sprintf("%s_%d", i, columns[j].ID)] {
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
	price, _ := strconv.Atoi(r.Form.Get("price"))

	db := connector.ConnectDB()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	column := schema.Column{Board: board.ID, Name: name, Price: price}

	db.Create(&column)
	db.Close()
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	columns := []schema.Column{}
	db.Where("board = ?", board.ID).Find(&columns)

	columnIdx, _ := strconv.Atoi(chi.URLParam(r, "column"))
	column := columns[columnIdx]

	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", chi.URLParam(r, "date")))
	t = t.Add(time.Hour * -9)
	check := schema.Check{Date: t, Column: column.ID, Board: board.ID}

	db.Create(&check)
	db.Close()
}


func uncheckHandler(w http.ResponseWriter, r *http.Request) {
	db := connector.ConnectDB()

	board := schema.Board{}
	db.First(&board, "token = ?", chi.URLParam(r, "boardToken"))

	columns := []schema.Column{}
	db.Where("board = ?", board.ID).Find(&columns)

	columnIdx, _ := strconv.Atoi(chi.URLParam(r, "column"))
	column := columns[columnIdx]

	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", chi.URLParam(r, "date")))
	t = t.Add(time.Hour * -9)

	db.Where("board = ? AND `column` = ? AND date = ?", board.ID, column.ID, t).Delete(&schema.Check{})
	db.Close()
}

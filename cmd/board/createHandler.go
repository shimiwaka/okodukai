package main

import (
	"fmt"
	"net/http"
	"crypto/md5"
	"time"

	// "github.com/go-chi/chi"
	// "github.com/jinzhu/gorm"
	"github.com/shimiwaka/okodukai/schema"
	"github.com/shimiwaka/okodukai/connector"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseForm()
	if e != nil {
		panic("error: parse error occured.")
	}

	email := r.Form.Get("email")

	db := connector.ConnectDB()

	seed := []byte(email + fmt.Sprint(time.Now().UnixNano()))
	token := fmt.Sprintf("%x", md5.Sum(seed))

	board := schema.Board{Owner: email, Token: token}
	db.Create(&board)

	db.Close()
	
	fmt.Fprintf(w, "{\"success\":\"true\", \"token\",\"%s\"}", token)
}

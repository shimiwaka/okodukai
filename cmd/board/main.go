package main

import (
	// "net/http"
	"net/http/cgi"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	rootPath := os.Getenv("SCRIPT_NAME")

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://pb.peraimaru.work"}}))
	
	r.Get(rootPath + "/", pingHandler)
	r.Post(rootPath + "/create", createHandler)
	r.Post(rootPath + "/forget", forgetHandler)

	r.Route(rootPath + "/board", func(r chi.Router) {
		r.Get("/{boardToken}", showBoardHandler)
		r.Post("/{boardToken}/newcolumn", addColumnHandler)
		r.Get("/{boardToken}/deletecolumn/{idx}", deleteColumnHandler)
		r.Get("/{boardToken}/check/{date}/{column}", checkHandler)
		r.Get("/{boardToken}/uncheck/{date}/{column}", uncheckHandler)
		r.Get("/{boardToken}/newpayment/{date}", newPaymentHandler)
		r.Get("/{boardToken}/cancelpayment/{date}", cancelPaymentHandler)
	  })

	// http.ListenAndServe(":9999", r)
	cgi.Serve(r)
}

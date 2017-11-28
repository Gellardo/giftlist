package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/Gellardo/giftlist/api"
	"github.com/Gellardo/giftlist/web"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	api.Setup(r, "/list/")
	web.Run(r, "/web/", "./web/", "http://localhost:8000/list/")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
	})

	server := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

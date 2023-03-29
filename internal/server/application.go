package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", SearchHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

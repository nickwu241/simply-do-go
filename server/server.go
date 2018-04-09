package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nickwu241/simply-do/server/store"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
)

type Server interface {
	Run(addr ...string)
}

func NewServer() (Server, error) {
	store, err := store.NewFirebaseStore()
	if err != nil {
		return nil, errors.Wrapf(err, "initializing firebase database")
	}
	api := API{
		store: store,
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/items", api.getItems).Methods("GET")
	router.HandleFunc("/api/items/{id}", api.getItem).Methods("GET")
	router.HandleFunc("/api/items", api.createItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", api.updateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", api.deleteItem).Methods("DELETE")

	CORSMiddleware := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
	server := negroni.Classic().With(negroni.HandlerFunc(CORSMiddleware))
	server.UseHandler(router)
	return server, nil
}
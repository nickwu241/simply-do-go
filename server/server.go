package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nickwu241/simply-do/server/store"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
)

// Server runs the server with the given addr or uses :8080 by default.
type Server interface {
	Run(addr ...string)
}

// NewServer returns an instance of Server.
func NewServer() (Server, error) {
	store, err := store.NewFirebaseStore()
	if err != nil {
		return nil, errors.Wrapf(err, "initializing firebase database")
	}
	api := API{store: store}
	router := mux.NewRouter()
	router.HandleFunc("/api/items", api.getItems).Methods("GET")
	router.HandleFunc("/api/items/{id}", api.getItem).Methods("GET")
	router.HandleFunc("/api/items", api.createItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", api.updateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", api.deleteItem).Methods("DELETE")
	defaultNotFoundHandler := router.NotFoundHandler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix("/#/", r.URL.Path) {
			http.Redirect(w, r, "/#/"+r.URL.Path, http.StatusPermanentRedirect)
		} else {
			defaultNotFoundHandler.ServeHTTP(w, r)
		}
	})

	CORSMiddleware := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-simply-do-uid")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		next(w, r)
	}
	server := negroni.New().With(
		negroni.HandlerFunc(CORSMiddleware),
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
		negroni.HandlerFunc(api.setUserMiddleware),
	)
	server.UseHandler(router)
	return server, nil
}

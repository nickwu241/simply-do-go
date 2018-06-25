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
	router := mux.NewRouter()
	api := API{store: store}
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Handle("/list/{lid}/items", api.setList(http.HandlerFunc(api.getItems))).Methods("GET")
	apiRouter.Handle("/list/{lid}/items/{id}", api.setList(http.HandlerFunc(api.getItem))).Methods("GET")
	apiRouter.Handle("/list/{lid}/items", api.setList(http.HandlerFunc(api.createItem))).Methods("POST")
	apiRouter.Handle("/list/{lid}/items/{id}", api.setList(http.HandlerFunc(api.updateItem))).Methods("PUT")
	apiRouter.Handle("/list/{lid}/items/{id}", api.setList(http.HandlerFunc(api.deleteItem))).Methods("DELETE")

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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		next(w, r)
	}
	server := negroni.Classic().With(
		negroni.HandlerFunc(CORSMiddleware),
	)
	server.UseHandler(router)
	return server, nil
}

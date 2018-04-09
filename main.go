package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type API struct {
	store Store
}

func main() {
	store, err := NewFirebaseStore()
	if err != nil {
		fmt.Printf("error initializing firebase database: %v\n", err)
		os.Exit(2)
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

	server := negroni.Classic().With(negroni.HandlerFunc(CORSMiddleware))
	server.UseHandler(router)

	// Use PORT from environment variables if it's set. Needed for Heroku.
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			fmt.Printf("error converting PORT environment to integer: %v\n", err)
			os.Exit(2)
		}
		server.Run(fmt.Sprintf(":%d", port))
	} else {
		server.Run(":8080")
	}
}

func CORSMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	next(w, r)
}

func (api *API) getItems(w http.ResponseWriter, r *http.Request) {
	items := api.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (api *API) getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item := api.store.Get(params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (api *API) createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		fmt.Printf("invalid request body: %v\n", err)
		return
	}
	item := api.store.Create(newItem)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (api *API) updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updateItem Item
	if err := json.NewDecoder(r.Body).Decode(&updateItem); err != nil {
		fmt.Printf("invalid request body: %v\n", err)
		return
	}
	item := api.store.Update(params["id"], updateItem)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (api *API) deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	items := api.store.Delete(params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

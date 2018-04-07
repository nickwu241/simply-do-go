package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var store = NewFirebaseStore()

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/items", getItems).Methods("GET")
	router.HandleFunc("/api/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/api/items", createItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", deleteItem).Methods("DELETE")

	server := negroni.Classic().With(negroni.HandlerFunc(CORSMiddleware))
	server.UseHandler(router)
	server.Run(":8080")
}

func CORSMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	next(w, r)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	items := store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item := store.Get(params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		fmt.Printf("invalid request body: %v\n", err)
		return
	}
	item := store.Create(newItem)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updateItem Item
	if err := json.NewDecoder(r.Body).Decode(&updateItem); err != nil {
		fmt.Printf("invalid request body: %v\n", err)
		return
	}
	item := store.Update(params["id"], updateItem)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	items := store.Delete(params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

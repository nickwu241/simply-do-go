package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nickwu241/simply-do/server/models"
	"github.com/nickwu241/simply-do/server/store"
)

// API contains the handlers for the server.
type API struct {
	store store.Store
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
	var newItem models.Item
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
	var updateItem models.Item
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

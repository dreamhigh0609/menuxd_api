package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/waiter"
)

// WaiterRouter is a router of the waiters.
type WaiterRouter struct {
	storage waiter.Storage
}

// getAllHandler response all the waiters from a client
func (wr WaiterRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	waiters, err := wr.storage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(waiters)
	if err != nil {
		http.Error(w, "Failed to parse waiters", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one waiter by id
func (wr WaiterRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wt, err := wr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(wt)
	if err != nil {
		http.Error(w, "Failed to parse waiters", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new waiter
func (wr WaiterRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	wt := waiter.New()
	err := json.NewDecoder(r.Body).Decode(wt)
	if err != nil {
		http.Error(w, "Invalid waiter", http.StatusBadRequest)
		return
	}

	err = wr.storage.Create(wt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored waiter by id
func (wr WaiterRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	wt := waiter.New()

	err := json.NewDecoder(r.Body).Decode(wt)
	if err != nil {
		http.Error(w, "Invalid waiter", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = wr.storage.Update(uint(id), wt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a waiter by ID
func (wr WaiterRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = wr.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewWaiterRouter inicialize a new router with each endpoint
func NewWaiterRouter(s waiter.Storage) *chi.Mux {
	r := chi.NewRouter()
	wr := WaiterRouter{storage: s}

	// Set endpoints
	r.Get("/client/{clientId}", wr.getAllHandler)
	r.Post("/", wr.createHandler)
	r.Get("/{id}", wr.getOneHandler)
	r.Put("/{id}", wr.updateHandler)
	r.Delete("/{id}", wr.deleteHandler)

	return r
}

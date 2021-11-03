package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/table"
)

// TableRouter is a router of the tables.
type TableRouter struct {
	storage table.Storage
}

// getAllHandler response all the tables from a client.
func (tr TableRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tables, err := tr.storage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(tables)
	if err != nil {
		http.Error(w, "Failed to parse tables", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one table by id.
func (tr TableRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := tr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Failed to parse tables", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new table.
func (tr TableRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	t := &table.Table{}
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		http.Error(w, "Invalid table", http.StatusBadRequest)
		return
	}

	err = tr.storage.Create(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored table by id.
func (tr TableRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	t := &table.Table{}

	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		http.Error(w, "Invalid table", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tr.storage.Update(uint(id), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a table by ID.
func (tr TableRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tr.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewTableRouter inicialize a new router with each endpoint.
func NewTableRouter(s table.Storage) *chi.Mux {
	r := chi.NewRouter()
	tr := TableRouter{storage: s}

	// Set endpoints
	r.Get("/client/{clientId}", tr.getAllHandler)
	r.Post("/", tr.createHandler)
	r.Get("/{id}", tr.getOneHandler)
	r.Put("/{id}", tr.updateHandler)
	r.Delete("/{id}", tr.deleteHandler)

	return r
}

package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/stay"
)

// StayRouter is a router of the stay.
type StayRouter struct {
	storage stay.Storage
}

// getAllHandler response all the stay from a client.
func (sr StayRouter) getByClientHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := sr.storage.GetByClient(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to parse stay", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getAllHandler response all the stay from a client.
func (sr StayRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	result, err := sr.storage.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to parse stay", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one stay by id.
func (sr StayRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	st, err := sr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(st)
	if err != nil {
		http.Error(w, "Failed to parse stay", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new stay.
func (sr StayRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	st := &stay.Stay{}
	err := json.NewDecoder(r.Body).Decode(st)
	if err != nil {
		http.Error(w, "Invalid stay", http.StatusBadRequest)
		return
	}

	err = sr.storage.Create(st)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewStayRouter inicialize a new router with each endpoint.
func NewStayRouter(s stay.Storage) *chi.Mux {
	r := chi.NewRouter()
	sr := StayRouter{storage: s}

	// Set endpoints
	r.Get("/", sr.getAllHandler)
	r.Get("/client/{clientId}", sr.getByClientHandler)
	r.Post("/", sr.createHandler)
	r.Get("/{id}", sr.getOneHandler)

	return r
}

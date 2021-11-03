package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/ad"
)

// AdRouter is a router to the ads.
type AdRouter struct {
	storage ad.Storage
}

// getAllHandler response all the ads from a client.
func (ar AdRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ads, err := ar.storage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(ads)
	if err != nil {
		http.Error(w, "Failed to parse ads", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one ad by id.
func (ar AdRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ad, err := ar.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(ad)
	if err != nil {
		http.Error(w, "Failed to parse ads", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new ad.
func (ar AdRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	a := &ad.Ad{}
	err := json.NewDecoder(r.Body).Decode(a)
	if err != nil {
		http.Error(w, "Failed to parse ad", http.StatusBadRequest)
		return
	}

	err = ar.storage.Create(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored ad by id.
func (ar AdRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	a := &ad.Ad{}
	err := json.NewDecoder(r.Body).Decode(a)
	if err != nil {
		http.Error(w, "Failed to parse ad", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ar.storage.Update(uint(id), a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// patchHandler update a stored ad by id.
func (ar AdRouter) patchHandler(w http.ResponseWriter, r *http.Request) {
	a := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Failed to parse ad", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ar.storage.Patch(uint(id), a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// addClickHandler add a new click to stored ad by id.
func (ar AdRouter) addClickHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ar.storage.AddClick(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler remove an ad by id.
func (ar AdRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ar.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewAdRouter inicialize a new router with each endpoint.
func NewAdRouter(s ad.Storage) *chi.Mux {
	r := chi.NewRouter()
	ar := AdRouter{storage: s}

	// Set endpoints.
	r.Get("/client/{clientId}", ar.getAllHandler)
	r.Post("/", ar.createHandler)
	r.Post("/add-click/{id}", ar.addClickHandler)
	r.Get("/{id}", ar.getOneHandler)
	r.Put("/{id}", ar.updateHandler)
	r.Patch("/{id}", ar.patchHandler)
	r.Delete("/{id}", ar.deleteHandler)

	return r
}

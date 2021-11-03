package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/picture"
)

// createHandler Create and returns a new picture
func createHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	clientID := chi.URLParam(r, "clientId")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	fileName := chi.URLParam(r, "name")
	if fileName == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	p, err := picture.Upload(r, fileName, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Failed to parse picture", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// NewPictureRouter inicialize a new user router with each endpoint
func NewPictureRouter() *chi.Mux {
	r := chi.NewRouter()

	// Set endpoints
	r.Post("/{clientId}/{name}", createHandler)

	return r
}

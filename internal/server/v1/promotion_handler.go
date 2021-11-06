package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/promotion"
)

// PromotionRouter is a router of the promotions.
type PromotionRouter struct {
	storage promotion.Storage
}

// getAllHandler response all the promotions from a client.
func (pr PromotionRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var promotions []promotion.Promotion
	allPromotions := r.URL.Query().Get("all")
	if allPromotions != "" {
		promotions, err = pr.storage.GetAll(uint(clientID))
	} else {
		promotions, err = pr.storage.GetAllActive(uint(clientID))
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(promotions)
	if err != nil {
		http.Error(w, "Failed to parse promotions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// addClickHandler add a new click to stored ad by id.
func (pr PromotionRouter) addClickHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pr.storage.AddClick(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// getOneHandler response one promotion by id.
func (pr PromotionRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := pr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Failed to parse promotions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new promotion.
func (pr PromotionRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	p := &promotion.Promotion{}
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		http.Error(w, "Failed to parse the promotion", http.StatusBadRequest)
		return
	}

	err = pr.storage.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored promotion by id.
func (pr PromotionRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	p := &promotion.Promotion{}
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		http.Error(w, "Failed to parse the promotion", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pr.storage.Update(uint(id), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a bill by ID.
func (pr PromotionRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pr.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewPromotionRouter inicialize a new router with each endpoint.
func NewPromotionRouter(s promotion.Storage) *chi.Mux {
	r := chi.NewRouter()
	pr := PromotionRouter{storage: s}

	// Set endpoints
	r.Get("/client/{clientId}", pr.getAllHandler)
	r.Post("/", pr.createHandler)
	r.Post("/add-click/{id}", pr.addClickHandler)
	r.Get("/{id}", pr.getOneHandler)
	r.Put("/{id}", pr.updateHandler)
	r.Delete("/{id}", pr.deleteHandler)

	return r
}

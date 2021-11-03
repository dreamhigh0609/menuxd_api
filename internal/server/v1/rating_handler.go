package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/rating"
)

// RatingRouter is the router of bills.
type RatingRouter struct {
	storage rating.Storage
}

// getAllHandler response all the bills from a client.
func (rr RatingRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	ratings, err := rr.storage.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, "Failed to parse ratings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getAllHandler response all the bills from a client.
func (rr RatingRouter) getAllByQuestionHandler(w http.ResponseWriter, r *http.Request) {
	questionIDStr := chi.URLParam(r, "questionId")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ratings, err := rr.storage.GetAllByQuestion(uint(questionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, "Failed to parse ratings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one rating by id.
func (rr RatingRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ra, err := rr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(ra)
	if err != nil {
		http.Error(w, "Failed to parse rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new rating.
func (rr RatingRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	ra := rating.Rating{}
	err := json.NewDecoder(r.Body).Decode(&ra)
	if err != nil {
		http.Error(w, "Failed to parse the rating", http.StatusBadRequest)
		return
	}

	err = rr.storage.Create(&ra)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewRatingRouter inicialize a new router with each endpoint.
func NewRatingRouter(s rating.Storage) *chi.Mux {
	r := chi.NewRouter()
	rr := RatingRouter{storage: s}

	// Set endpoints
	r.Get("/", rr.getAllHandler)
	r.Get("/question/{questionId}", rr.getAllByQuestionHandler)
	r.Post("/", rr.createHandler)
	r.Get("/{id}", rr.getOneHandler)

	return r
}

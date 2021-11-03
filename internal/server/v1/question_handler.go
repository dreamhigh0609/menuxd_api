package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/question"
)

// QuestionRouter is a router of the questions.
type QuestionRouter struct {
	storage question.Storage
}

// getAllHandler response all the questions from a client.
func (qr QuestionRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	questions, err := qr.storage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(questions)
	if err != nil {
		http.Error(w, "Failed to parse questions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one question by id.
func (qr QuestionRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q, err := qr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(q)
	if err != nil {
		http.Error(w, "Failed to parse questions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new question.
func (qr QuestionRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	q := &question.Question{}
	err := json.NewDecoder(r.Body).Decode(q)
	if err != nil {
		http.Error(w, "Invalid question", http.StatusBadRequest)
		return
	}

	err = qr.storage.Create(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored question by id.
func (qr QuestionRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	q := &question.Question{}

	err := json.NewDecoder(r.Body).Decode(q)
	if err != nil {
		http.Error(w, "Invalid question", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = qr.storage.Update(uint(id), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a question by ID.
func (qr QuestionRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = qr.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewQuestionRouter inicialize a new router with each endpoint.
func NewQuestionRouter(s question.Storage) *chi.Mux {
	r := chi.NewRouter()
	qr := QuestionRouter{storage: s}

	// Set endpoints
	r.Get("/client/{clientId}", qr.getAllHandler)
	r.Post("/", qr.createHandler)
	r.Get("/{id}", qr.getOneHandler)
	r.Put("/{id}", qr.updateHandler)
	r.Delete("/{id}", qr.deleteHandler)

	return r
}

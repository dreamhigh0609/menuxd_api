package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/super-radmir/menuxd_api/pkg/category"
)

// CategoryRouter is a router to Categories.
type CategoryRouter struct {
	storage category.Storage
}

// getAllHandler response all the categories from a client.
func (cr CategoryRouter) getAllActiveHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := cr.storage.GetAllActive(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Failed to parse categories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getAllHandler response all the categories from a client.
func (cr CategoryRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := cr.storage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Failed to parse categories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getAllByBackupHandler response all the categories from a client.
func (cr CategoryRouter) getAllByBackupHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dishes, err := cr.storage.GetAllBackup(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(dishes)
	if err != nil {
		http.Error(w, "Failed to parse dishes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Disposition", "Attachment")
	http.ServeContent(w, r, "categories.json", time.Now(), bytes.NewReader(j))
}

// getOneHandler response one category by id.
func (cr CategoryRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := cr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(c)
	if err != nil {
		http.Error(w, "Failed to parse categories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new category.
func (cr CategoryRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	c := &category.Category{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		http.Error(w, "Failed to parse category", http.StatusBadRequest)
		return
	}

	err = cr.storage.Create(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// createManyHandler Create many categories from a file.
func (cr CategoryRouter) createManyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("categories")
	if err != nil {
		http.Error(w, "Failed to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categories := []category.Category{}
	err = json.NewDecoder(file).Decode(&categories)
	if err != nil {
		http.Error(w, "Failed to parse dishes", http.StatusBadRequest)
		return
	}

	err = cr.storage.CreateMany(uint(clientID), categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// createManyHandler Create many categories from a file.
func (cr CategoryRouter) updatePositionHandler(w http.ResponseWriter, r *http.Request) {
	categories := []category.Category{}
	err := json.NewDecoder(r.Body).Decode(&categories)
	if err != nil {
		http.Error(w, "Failed to parse dishes", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = cr.storage.UpdatePositions(categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored category by id.
func (cr CategoryRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	c := category.New()
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		http.Error(w, "Failed to parse category", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cr.storage.Update(uint(id), c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// patchHandler update a part of the stored category by id.
func (cr CategoryRouter) patchHandler(w http.ResponseWriter, r *http.Request) {
	c := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Failed to parse category", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cr.storage.Patch(uint(id), c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a category by ID.
func (cr CategoryRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cr.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewCategoryRouter inicialize a new router with each endpoint.
func NewCategoryRouter(s category.Storage) *chi.Mux {
	r := chi.NewRouter()
	cr := CategoryRouter{storage: s}
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)

	// Set endpoints.
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/client/{clientId}", cr.getAllActiveHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/client/{clientId}/admin", cr.getAllHandler)
	r.Get("/client/{clientId}/categories.json", cr.getAllByBackupHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Post("/client/{clientId}", cr.createManyHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Post("/", cr.createHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/{id}", cr.getOneHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Put("/{id}", cr.updateHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Delete("/{id}", cr.deleteHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Patch("/{id}", cr.patchHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Put("/position/", cr.updatePositionHandler)

	return r
}

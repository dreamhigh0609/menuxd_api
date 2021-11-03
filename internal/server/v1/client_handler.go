package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/super-radmir/menuxd_api/pkg/client"
	"github.com/super-radmir/menuxd_api/pkg/middleware/auth"
)

// ClientRouter is the router of clients.
type ClientRouter struct {
	storage client.Storage
}

// getAllHandler response all the clients.
func (cr ClientRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clients, err := cr.storage.GetAll(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, "Failed to parse clients", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one client by id.
func (cr ClientRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to parse clients", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new client.
func (cr ClientRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	c := client.New()

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Failed to parse client", http.StatusBadRequest)
		return
	}

	err = cr.storage.Create(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored client by id.
func (cr ClientRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	c := client.New()

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Failed to parse client", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cr.storage.Update(uint(id), &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a client by ID.
func (cr ClientRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
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

// NewClientRouter inicialize a new client router with each endpoint.
func NewClientRouter(s client.Storage) *chi.Mux {
	r := chi.NewRouter()
	cr := ClientRouter{storage: s}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)

	// Set endpoints.
	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Post("/", cr.createHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/{id}", cr.getOneHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Put("/{id}", cr.updateHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Delete("/{id}", cr.deleteHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/user/{userId}", cr.getAllHandler)

	return r
}

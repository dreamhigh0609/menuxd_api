package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/super-radmir/menuxd_api/pkg/middleware/auth"
	"github.com/super-radmir/menuxd_api/pkg/user"
)

// UserRouter is the router of the users.
type UserRouter struct {
	storage user.Storage
}

// getAllHandler response all the users
func (ur UserRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := ur.storage.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to parse users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one user by id.
func (ur UserRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := ur.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Failed to parse users", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new user.
func (ur UserRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	u := user.New()

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, "Failed to parse user", http.StatusBadRequest)
		return
	}

	err = ur.storage.Create(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored user by id.
func (ur UserRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	u := user.New()

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, "Failed to parse user", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ur.storage.Update(uint(id), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a user by ID.
func (ur UserRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ur.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// LoginHandler handle login, returns a json with token.
func (ur UserRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := user.New()
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, "Failed to parse user", http.StatusBadRequest)
		return
	}

	storedUser, err := ur.storage.GetByEmail(u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !storedUser.ComparePass(u.Password) {
		http.Error(w, "Passwords don't match", http.StatusUnauthorized)
		return
	}

	storedUser.CleanPass()
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)
	claims := jwtauth.Claims{
		"id":   strconv.Itoa(int(storedUser.ID)),
		"role": storedUser.Role,
	}

	_, token, err := tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	data := map[string]interface{}{
		"token": token,
		"user":  storedUser,
	}

	j, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// confirmUserHandler update user, set confirm to true.
func (ur UserRouter) confirmUserHandler(w http.ResponseWriter, r *http.Request) {
	u := user.New()
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, "Failed to parse user", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedUser, err := ur.storage.GetByEmail(u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !storedUser.ComparePass(u.OldPassword) {
		http.Error(w, "Passwords don't match", http.StatusUnauthorized)
		return
	}

	err = ur.storage.Confirm(uint(id), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// ForgotPasswordHandler update user, set confirm to true.
func (ur UserRouter) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	u := user.New()
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, "Failed to parse user", http.StatusBadRequest)
		return
	}

	err = ur.storage.RecoverPassword(u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewUserRouter inicialize a new user router with each endpoint.
func NewUserRouter(s user.Storage) (*chi.Mux, UserRouter) {
	r := chi.NewRouter()
	ur := UserRouter{storage: s}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)

	// Set endpoints
	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Put("/change-password/{id}", ur.confirmUserHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Get("/", ur.getAllHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Post("/", ur.createHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/{id}", ur.getOneHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Put("/{id}", ur.updateHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("admin"),
	).Delete("/{id}", ur.deleteHandler)

	return r, ur
}

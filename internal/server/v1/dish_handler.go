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
	"github.com/super-radmir/menuxd_api/pkg/dish"
	"github.com/super-radmir/menuxd_api/pkg/middleware/auth"
)

// DishRouter is a router to dishes.
type DishRouter struct {
	storage dish.Storage
}

// getAllPaginateHandler response all the dishes from a client.
func (dr DishRouter) getAllPaginateHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageStr := chi.URLParam(r, "page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Failed to parse page", http.StatusBadRequest)
		return
	}

	dishes, total, err := dr.storage.GetAllWithPagination(uint(clientID), int64(page))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	result := map[string]interface{}{
		"dishes": dishes,
		"total":  total,
	}

	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to parse dishes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getAllByBackupHandler response all the dishes from a client.
func (dr DishRouter) getAllByBackupHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dishes, err := dr.storage.GetAllBackup(uint(clientID))
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
	http.ServeContent(w, r, "dishes.json", time.Now(), bytes.NewReader(j))
}

// getAllByBackupHandler response all the dishes from a client.
func (dr DishRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dishes, err := dr.storage.GetAll(uint(clientID))
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getByCategory return a list of dishes by category.
func (dr DishRouter) getByCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "categoryId")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allDishes := r.URL.Query().Get("all")
	var dishes []dish.Dish
	if allDishes != "" {
		dishes, err = dr.storage.GetAllByCategory(uint(categoryID))
	} else {
		dishes, err = dr.storage.GetAllActiveByCategory(uint(categoryID))
	}
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one dish by id.
func (dr DishRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := dr.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Failed to parse dishes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new dish.
func (dr DishRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	d := &dish.Dish{}
	err := json.NewDecoder(r.Body).Decode(d)
	if err != nil {
		http.Error(w, "Failed to parse dish", http.StatusBadRequest)
		return
	}

	err = dr.storage.Create(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// createManyHandler Create many dishes by json file.
func (dr DishRouter) createManyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("dishes")
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

	var dishes dish.Dishes
	err = json.NewDecoder(file).Decode(&dishes)
	if err != nil {
		http.Error(w, "Failed to parse dishes", http.StatusBadRequest)
		return
	}

	err = dr.storage.CreateMany(uint(clientID), dishes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored dish by id.
func (dr DishRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	d := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Failed to parse dish", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dr.storage.Update(uint(id), d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a dish by ID.
func (dr DishRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dr.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// getSuggestedHandler Remove a dish by ID.
func (dr DishRouter) getSuggestedHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "categoryId")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dishes, err := dr.storage.GetSuggested(uint(categoryID))
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// addClickHandler add a new click to suggested dish by id.
func (dr DishRouter) addClickHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dr.storage.AddClick(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// getClicksHandler get clicks by client id.
func (dr DishRouter) getClicksHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clicks, err := dr.storage.GetClicks(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(clicks)
	if err != nil {
		http.Error(w, "Failed to parse clicks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// NewDishRouter inicialize a new router with each endpoint.
func NewDishRouter(s dish.Storage) *chi.Mux {
	r := chi.NewRouter()
	dr := DishRouter{storage: s}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)
	jwtauth.Verifier(tokenAuth)
	// Set endpoints.
	r.Get("/client/{clientId}/dishes.json", dr.getAllByBackupHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/client/{clientId}", dr.getAllHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/client/{clientId}/{page:[0-9]+}", dr.getAllPaginateHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(auth.Authenticator("client")).Post("/client/{clientId}", dr.createManyHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/category/{categoryId}", dr.getByCategory)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Post("/", dr.createHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/{id}", dr.getOneHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Put("/{id}", dr.updateHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Delete("/{id}", dr.deleteHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Get("/suggested/{categoryId}", dr.getSuggestedHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Post("/add-click/{id}", dr.addClickHandler)

	r.With(
		jwtauth.Verifier(tokenAuth),
	).With(
		auth.Authenticator("client"),
	).Post("/clicks/client/{clientId}", dr.getClicksHandler)

	return r
}

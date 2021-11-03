package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/super-radmir/menuxd_api/pkg/bill"
)

// BillRouter is the router of bills.
type BillRouter struct {
	storage bill.Storage
}

// getAllHandler response all the bills from a client.
func (br BillRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bills, err := br.storage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(bills)
	if err != nil {
		http.Error(w, "Failed to parse bills", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// getOneHandler response one bill by id.
func (br BillRouter) getOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := br.storage.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(b)
	if err != nil {
		http.Error(w, "Failed to parse bills", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// createHandler Create a new bill.
func (br BillRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	b := bill.New()
	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, "Failed to parse the bill", http.StatusBadRequest)
		return
	}

	err = br.storage.Create(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// updateHandler update a stored bill by id.
func (br BillRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	b := bill.New()

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, "Failed to parse the bill", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = br.storage.Update(uint(id), b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// deleteHandler Remove a bill by ID.
func (br BillRouter) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = br.storage.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// NewBillRouter inicialize a new router with each endpoint.
func NewBillRouter(s bill.Storage) *chi.Mux {
	r := chi.NewRouter()
	br := BillRouter{storage: s}

	// Set endpoints
	r.Get("/client/{clientId}", br.getAllHandler)
	r.Post("/", br.createHandler)
	r.Get("/{id}", br.getOneHandler)
	r.Put("/{id}", br.updateHandler)
	r.Delete("/{id}", br.deleteHandler)

	return r
}

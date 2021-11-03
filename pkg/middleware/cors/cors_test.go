package cors

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestCorsAllAllowed(t *testing.T) {
	mux := chi.NewRouter()
	mux.Use(Cors())
	mux.Get("/cors", fakeHanlder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/cors", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected code %d, got code %d", http.StatusOK, w.Code)
	}
	allow := w.Header().Get("Access-Control-Allow-Origin")
	if allow != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin = *, got %s", allow)
	}
}

func TestCorsWhiteList(t *testing.T) {
	mux := chi.NewRouter()
	allowedList := []string{"https://example.com", "https://anotherexample.com"}

	mux.Use(Cors(allowedList...))
	mux.Get("/cors", fakeHanlder)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/cors", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected code %d, got code %d", http.StatusOK, w.Code)
	}
	allow := w.Header().Get("Access-Control-Allow-Origin")
	if allow != strings.Join(allowedList, ", ") {
		t.Errorf("Expected Access-Control-Allow-Origin = %s, got %s", strings.Join(allowedList, ", "), allow)
	}
}

func fakeHanlder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world"))
}

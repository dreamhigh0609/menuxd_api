package cors

import (
	"net/http"
	"strings"
)

// Cors set the header Access-Control-Allow-Origin to gave hosts, if empty all host allowed
func Cors(url ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		length := len(url)
		var allowed string
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if length == 0 {
				allowed = "*"
			} else {
				allowed = strings.Join(url, ", ")
			}
			w.Header().Set("Access-Control-Allow-Origin", allowed)
			next.ServeHTTP(w, r)
		})
	}
}

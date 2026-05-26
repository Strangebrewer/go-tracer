package middleware

import (
	"net/http"
	"strings"
)

func StripAPIPrefix(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, prefix) {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
				r.URL.RawPath = ""
			}
			next.ServeHTTP(w, r)
		})
	}
}

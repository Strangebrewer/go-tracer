package span

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(store *Store, serviceKeyMiddleware, authMiddleware func(http.Handler) http.Handler) chi.Router {
	h := newHandler(store)
	r := chi.NewRouter()

	r.With(serviceKeyMiddleware).Post("/spans", h.createSpan)
	r.With(authMiddleware).Get("/traces/{traceId}", h.getTrace)

	return r
}

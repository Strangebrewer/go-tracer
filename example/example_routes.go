package example

import (
	"github.com/go-chi/chi/v5"
)

func Routes(store *Store) chi.Router {
	h := NewHandler(store)
	r := chi.NewRouter()

	r.Get("/", h.GetAll)
	r.Get("/{id}", h.GetOne)
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)

	return r
}

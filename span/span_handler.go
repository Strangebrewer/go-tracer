package span

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	store *Store
}

func newHandler(store *Store) *handler {
	return &handler{store: store}
}

func (h *handler) createSpan(w http.ResponseWriter, r *http.Request) {
	var input CreateSpanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.store.Create(r.Context(), input); err != nil {
		slog.Error("failed to create span", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) getTrace(w http.ResponseWriter, r *http.Request) {
	traceID := chi.URLParam(r, "traceId")

	spans, err := h.store.GetByTraceID(r.Context(), traceID)
	if err != nil {
		slog.Error("failed to get trace", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spans)
}

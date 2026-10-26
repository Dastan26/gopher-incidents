// Package httpapi — HTTP-слой: разбирает запросы, зовёт хранилище,
// пишет ответы в формате JSON:API.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"gopher-incidents/internal/brand"
)

// Repository описан здесь, у потребителя, а не рядом с реализацией.
type Repository interface {
	List(ctx context.Context) ([]brand.Brand, error)
	Get(ctx context.Context, id int) (brand.Brand, error)
	Create(ctx context.Context, b brand.Brand) (brand.Brand, error)
	Update(ctx context.Context, b brand.Brand) error
	Delete(ctx context.Context, id int) error
}

type handler struct {
	repo Repository
	log  *slog.Logger
}

// New собирает роутер со всеми middleware.
func New(repo Repository, log *slog.Logger, allowedOrigin string) http.Handler {
	h := &handler{repo: repo, log: log}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /brands", h.list)
	mux.HandleFunc("GET /brands/{id}", h.get)
	mux.HandleFunc("POST /brands", h.create)
	mux.HandleFunc("PATCH /brands/{id}", h.update)
	mux.HandleFunc("DELETE /brands/{id}", h.delete)

	return withLogging(log, withCORS(allowedOrigin, mux))
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	out := make([]resource, 0, len(items))
	for _, b := range items {
		out = append(out, toResource(b))
	}
	writeJSON(w, http.StatusOK, document{Data: out})
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	b, err := h.repo.Get(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, document{Data: toResource(b)})
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var in requestDocument
	if !decode(w, r, &in) {
		return
	}
	attrs := in.Data.Attributes
	b := brand.Brand{NameKz: deref(attrs.NameKz), NameRu: deref(attrs.NameRu)}
	if err := b.Validate(); err != nil {
		h.fail(w, r, err)
		return
	}
	created, err := h.repo.Create(r.Context(), b)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, document{Data: toResource(created)})
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in requestDocument
	if !decode(w, r, &in) {
		return
	}
	b, err := h.repo.Get(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	attrs := in.Data.Attributes
	if attrs.NameKz != nil {
		b.NameKz = *attrs.NameKz
	}
	if attrs.NameRu != nil {
		b.NameRu = *attrs.NameRu
	}
	if err := b.Validate(); err != nil {
		h.fail(w, r, err)
		return
	}
	if err := h.repo.Update(r.Context(), b); err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, document{Data: toResource(b)})
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// fail превращает ошибку из нижних слоёв в ответ. Только здесь решается,
// какой HTTP-статус соответствует какой ошибке.
func (h *handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	var ve *brand.ValidationError
	switch {
	case errors.Is(err, brand.ErrNotFound):
		writeError(w, http.StatusNotFound, "Бренд не найден", "")
	case errors.As(err, &ve):
		writeError(w, http.StatusUnprocessableEntity, ve.Msg, "/data/attributes/"+ve.Field)
	default:
		h.log.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
		writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера", "")
	}
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusNotFound, "Бренд не найден", "")
		return 0, false
	}
	return id, true
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный JSON в теле запроса", "")
		return false
	}
	return true
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

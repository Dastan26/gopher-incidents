package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopher-incidents/internal/brand"
	"gopher-incidents/internal/store"
)

func newTestServer(t *testing.T) (http.Handler, *store.Memory) {
	t.Helper()
	repo := store.NewMemory()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(repo, log, "http://localhost:5173"), repo
}

func TestListBrands(t *testing.T) {
	srv, repo := newTestServer(t)
	if _, err := repo.Create(context.Background(), brand.Brand{NameKz: "Рахат", NameRu: "Рахат"}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/brands", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ждали 200", rec.Code)
	}
	var doc struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&doc); err != nil {
		t.Fatalf("тело не JSON: %v", err)
	}
	if len(doc.Data) != 1 || doc.Data[0].ID != "1" {
		t.Errorf("data = %+v, ждали один бренд с id 1", doc.Data)
	}
}

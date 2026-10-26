package store

import (
	"context"
	"testing"

	"gopher-incidents/internal/brand"
)

func TestMemoryCreateAndList(t *testing.T) {
	ctx := context.Background()
	m := NewMemory()

	for _, name := range []string{"Рахат", "Шоро"} {
		if _, err := m.Create(ctx, brand.Brand{NameRu: name, NameKz: name}); err != nil {
			t.Fatalf("Create(%q): %v", name, err)
		}
	}

	got, err := m.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Errorf("List = %+v, ждали два бренда с id 1 и 2", got)
	}
}

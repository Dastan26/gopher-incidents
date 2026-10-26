// Package store — хранилище брендов. Пока в памяти; на следующих уровнях
// рядом появится реализация на Postgres с тем же набором методов.
package store

import (
	"context"
	"sort"
	"sync"

	"gopher-incidents/internal/brand"
)

type Memory struct {
	mu     sync.RWMutex
	nextID int
	items  map[int]brand.Brand
}

func NewMemory() *Memory {
	return &Memory{nextID: 1, items: make(map[int]brand.Brand)}
}

func (m *Memory) List(_ context.Context) ([]brand.Brand, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]brand.Brand, 0, len(m.items))
	for _, b := range m.items {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *Memory) Get(_ context.Context, id int) (brand.Brand, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.items[id]
	if !ok {
		return brand.Brand{}, brand.ErrNotFound
	}
	return b, nil
}

func (m *Memory) Create(_ context.Context, b brand.Brand) (brand.Brand, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	b.ID = m.nextID
	m.nextID++
	m.items[b.ID] = b
	return b, nil
}

func (m *Memory) Update(_ context.Context, b brand.Brand) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.items[b.ID]; !ok {
		return brand.ErrNotFound
	}
	m.items[b.ID] = b
	return nil
}

func (m *Memory) Delete(_ context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.items[id]; !ok {
		return brand.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

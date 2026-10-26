// Package brand — предметная область: бренд и правила его валидации.
// Пакет ничего не знает ни про HTTP, ни про хранилище.
package brand

import "errors"

// MaxNameLen — максимальная длина названия бренда в символах (не в байтах).
const MaxNameLen = 100

// ErrNotFound возвращает хранилище, если бренда с таким id нет.
var ErrNotFound = errors.New("brand not found")

type Brand struct {
	ID     int
	NameRu string
	NameKz string
}

// ValidationError — ошибка в конкретном поле. Field совпадает с именем
// атрибута в JSON:API (nameRu, nameKz), чтобы хендлер собрал source.pointer.
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string { return e.Field + ": " + e.Msg }

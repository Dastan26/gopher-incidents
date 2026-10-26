package brand

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ValidateName проверяет одно название: обязательное, пробелы по краям
// не считаются, не длиннее MaxNameLen символов на любом языке.
func ValidateName(field, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return &ValidationError{Field: field, Msg: "Обязательное поле"}
	}
	if utf8.RuneCountInString(name) > MaxNameLen {
		return &ValidationError{Field: field, Msg: fmt.Sprintf("Не длиннее %d символов", MaxNameLen)}
	}
	return nil
}

// Validate проверяет бренд целиком. Казахское название идёт первым,
// как и в форме на фронте.
func (b Brand) Validate() error {
	if err := ValidateName("nameKz", b.NameKz); err != nil {
		return err
	}
	return ValidateName("nameRu", b.NameRu)
}

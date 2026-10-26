package brand

import (
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{"пустое", "", true},
		{"только пробелы", "   ", true},
		{"латиница", "Nike", false},
		{"кириллица", "Рахат", false},
		{"ровно 100 символов кириллицей", strings.Repeat("ж", 100), false},
		{"101 символ кириллицей", strings.Repeat("ж", 101), true},
		{"101 символ латиницей", strings.Repeat("a", 101), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName("nameRu", tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) err = %v, wantErr = %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

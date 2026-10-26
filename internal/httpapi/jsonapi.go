package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gopher-incidents/internal/brand"
)

// mediaType — Content-Type ответов JSON:API, как у mdm-pimcore.
const mediaType = "application/vnd.api+json"

const brandType = "brands"

type document struct {
	Data any `json:"data"`
}

type resource struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes brandAttrs `json:"attributes"`
	Meta       meta       `json:"meta"`
}

type brandAttrs struct {
	NameKz string `json:"nameKz"`
	NameRu string `json:"nameRu"`
}

type meta struct {
	Permissions permissions `json:"permissions"`
}

type permissions struct {
	Save   bool `json:"save"`
	Delete bool `json:"delete"`
}

// requestDocument — тело POST и PATCH. Атрибуты указателями: в PATCH
// клиент присылает только изменённые поля, и nil значит «не трогать».
type requestDocument struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			NameKz *string `json:"nameKz"`
			NameRu *string `json:"nameRu"`
		} `json:"attributes"`
	} `json:"data"`
}

type errorDocument struct {
	Errors []errorObject `json:"errors"`
}

type errorObject struct {
	Status string       `json:"status"`
	Detail string       `json:"detail"`
	Source *errorSource `json:"source,omitempty"`
}

type errorSource struct {
	Pointer string `json:"pointer"`
}

func toResource(b brand.Brand) resource {
	return resource{
		Type:       brandType,
		ID:         strconv.Itoa(b.ID),
		Attributes: brandAttrs{NameKz: b.NameKz, NameRu: b.NameRu},
		// Авторизации пока нет, поэтому всё разрешено. Появится на уровне про Keycloak.
		Meta: meta{Permissions: permissions{Save: true, Delete: true}},
	}
}

// writeJSON пишет ответ. Порядок важен: заголовки, затем статус, затем тело.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", mediaType)
	w.WriteHeader(status)
	// Статус уже ушёл клиенту, поменять его при ошибке кодирования нельзя.
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, detail, pointer string) {
	e := errorObject{Status: strconv.Itoa(status), Detail: detail}
	if pointer != "" {
		e.Source = &errorSource{Pointer: pointer}
	}
	writeJSON(w, status, errorDocument{Errors: []errorObject{e}})
}

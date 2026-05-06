package rest

import (
	"net/http"

	controllerv1 "go.yunus-emre.dev/url-shortaner/apiserver/internal/api/rest/v1"
	"go.yunus-emre.dev/url-shortaner/storage"
)

func RegisterRoutes(mux *http.ServeMux, storage storage.Storage) {
	controllerv1.New(storage).RegisterRoutes(mux)
}

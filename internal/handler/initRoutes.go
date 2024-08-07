package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes(h *URLHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/shorten", h.ShortenURL).Methods(http.MethodPost)
	r.HandleFunc("/{shortURL}", h.ResolveURL).Methods(http.MethodGet)

	return r

}
